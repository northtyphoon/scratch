package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand/v2"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

// Simple Redis cluster connection maintainer opening N concurrent connections.
// It uses raw TCP/TLS with lightweight PING keepalives and optional AUTH.
// Designed for load/soak testing connection capacity.

// Contract:
// - Inputs: redis addresses (comma-separated), desired connections (N), auth, db, useTLS, ping interval, connect rate, per-conn idle timeout, per-conn read timeout
// - Behavior: dials connections in a round-robin fashion across nodes, keeps them alive with periodic PING, reconnects on error, tracks metrics
// - Outputs: logs and periodic metrics; exit on SIGINT/SIGTERM

// Edge cases handled: slow connects, TLS, auth, cluster MOVED responses ignored (we don't issue keys), backoff on failures, rate-limited dial bursts.

type dialTarget struct {
	addr   string
	useTLS bool
}

func dialOne(ctx context.Context, t dialTarget, timeout time.Duration) (net.Conn, error) {
	dialer := &net.Dialer{Timeout: timeout, KeepAlive: 30 * time.Second}
	if !t.useTLS {
		return dialer.DialContext(ctx, "tcp", t.addr)
	}
	tlsCfg := &tls.Config{InsecureSkipVerify: true} // for testing; optionally pin certs
	return tls.DialWithDialer(dialer, "tcp", t.addr, tlsCfg)
}

func writeRESPPing(w *bufio.Writer) error {
	// *1\r\n$4\r\nPING\r\n
	if _, err := w.WriteString("*1\r\n$4\r\nPING\r\n"); err != nil {
		return err
	}
	return w.Flush()
}

func writeRESPAuth(w *bufio.Writer, username, password string) error {
	// AUTH <username> <password> or AUTH <password>
	if username != "" {
		line := fmt.Sprintf("*3\r\n$4\r\nAUTH\r\n$%d\r\n%s\r\n$%d\r\n%s\r\n", len(username), username, len(password), password)
		if _, err := w.WriteString(line); err != nil {
			return err
		}
		return w.Flush()
	}
	line := fmt.Sprintf("*2\r\n$4\r\nAUTH\r\n$%d\r\n%s\r\n", len(password), password)
	if _, err := w.WriteString(line); err != nil {
		return err
	}
	return w.Flush()
}

func writeRESPSelect(w *bufio.Writer, db int) error {
	line := fmt.Sprintf("*2\r\n$6\r\nSELECT\r\n$%d\r\n%d\r\n", len(strconv.Itoa(db)), db)
	if _, err := w.WriteString(line); err != nil {
		return err
	}
	return w.Flush()
}

func readSimpleRESP(r *bufio.Reader, deadline time.Time) error {
	if !deadline.IsZero() {
		_ = r.Buffered() // noop to appease lints; we set deadlines on conn externally
	}
	b, err := r.ReadByte()
	if err != nil {
		return err
	}
	switch b {
	case '+':
		// Simple string
		_, err = r.ReadString('\n')
		return err
	case '-':
		// Error
		line, _ := r.ReadString('\n')
		return errors.New(strings.TrimSpace(line))
	case '$':
		// Bulk string: read length, then payload and CRLF
		lenLine, err := r.ReadString('\n')
		if err != nil {
			return err
		}
		lenLine = strings.TrimSpace(lenLine)
		n, err := strconv.Atoi(lenLine)
		if err != nil {
			return err
		}
		if n >= 0 {
			buf := make([]byte, n+2)
			_, err = r.Read(buf)
			return err
		}
		return nil
	case '*':
		// Array: skip line
		_, err = r.ReadString('\n')
		return err
	default:
		// Unexpected
		rest, _ := r.ReadString('\n')
		return fmt.Errorf("unexpected RESP prefix %q rest %q", b, rest)
	}
}

type connWorkerCfg struct {
	id           int
	targets      []dialTarget
	username     string
	password     string
	db           int
	pingEvery    time.Duration
	readTimeout  time.Duration
	idleTimeout  time.Duration
	connectTO    time.Duration
	backoffMin   time.Duration
	backoffMax   time.Duration
}

type metrics struct {
	opened    atomic.Int64
	authed    atomic.Int64
	pongs     atomic.Int64
	errors    atomic.Int64
	restarts  atomic.Int64
	active    atomic.Int64
}

func jitter(d time.Duration, pct float64) time.Duration {
	if d <= 0 || pct <= 0 {
		return d
	}
	jd := time.Duration(float64(d) * (1 + (rand.Float64()*2-1)*pct))
	if jd < 0 {
		return 0
	}
	return jd
}

func backoffSeq(min, max time.Duration) func(int) time.Duration {
	if min <= 0 {
		min = 100 * time.Millisecond
	}
	if max < min {
		max = min * 10
	}
	return func(retry int) time.Duration {
		b := min << retry
		if b > max {
			b = max
		}
		return jitter(b, 0.2)
	}
}

func connWorker(ctx context.Context, cfg connWorkerCfg, m *metrics) {
	bseq := backoffSeq(cfg.backoffMin, cfg.backoffMax)
	nextTarget := cfg.id
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		nextTarget = (nextTarget + 1) % len(cfg.targets)
		t := cfg.targets[nextTarget]

		connCtx, cancel := context.WithTimeout(ctx, cfg.connectTO)
		c, err := dialOne(connCtx, t, cfg.connectTO)
		cancel()
		if err != nil {
			m.errors.Add(1)
			log.Printf("conn[%d] dial to %s failed: %v", cfg.id, t.addr, err)
			wait := bseq(int(m.restarts.Load() % 10))
			time.Sleep(wait)
			continue
		}
		m.opened.Add(1)
		m.active.Add(1)

		br := bufio.NewReader(c)
		bw := bufio.NewWriter(c)

		// Optional AUTH
		if cfg.password != "" || cfg.username != "" {
			if err := writeRESPAuth(bw, cfg.username, cfg.password); err != nil {
				m.errors.Add(1)
				log.Printf("conn[%d] AUTH write to %s failed: %v", cfg.id, t.addr, err)
				_ = c.Close()
				m.active.Add(^int64(0))
				continue
			}
			_ = c.SetReadDeadline(time.Now().Add(cfg.readTimeout))
			if err := readSimpleRESP(br, time.Now().Add(cfg.readTimeout)); err != nil {
				m.errors.Add(1)
				log.Printf("conn[%d] AUTH response from %s error: %v", cfg.id, t.addr, err)
				_ = c.Close()
				m.active.Add(^int64(0))
				continue
			}
			m.authed.Add(1)
		}

		// Optional SELECT DB
		if cfg.db > 0 {
			if err := writeRESPSelect(bw, cfg.db); err != nil {
				m.errors.Add(1)
				log.Printf("conn[%d] SELECT write to %s failed (db=%d): %v", cfg.id, t.addr, cfg.db, err)
				_ = c.Close()
				m.active.Add(^int64(0))
				continue
			}
			_ = c.SetReadDeadline(time.Now().Add(cfg.readTimeout))
			if err := readSimpleRESP(br, time.Now().Add(cfg.readTimeout)); err != nil {
				m.errors.Add(1)
				log.Printf("conn[%d] SELECT response from %s error: %v", cfg.id, t.addr, err)
				_ = c.Close()
				m.active.Add(^int64(0))
				continue
			}
		}

		lastIO := time.Now()
		pingTicker := time.NewTicker(jitter(cfg.pingEvery, 0.1))
		idleTimer := time.NewTimer(cfg.idleTimeout)

		alive := true
		for alive {
			select {
			case <-ctx.Done():
				alive = false
			case <-pingTicker.C:
				if err := writeRESPPing(bw); err != nil {
					m.errors.Add(1)
					log.Printf("conn[%d] PING write to %s failed: %v", cfg.id, t.addr, err)
					alive = false
					break
				}
				_ = c.SetReadDeadline(time.Now().Add(cfg.readTimeout))
				if err := readSimpleRESP(br, time.Now().Add(cfg.readTimeout)); err != nil {
					m.errors.Add(1)
					log.Printf("conn[%d] PING response from %s error: %v", cfg.id, t.addr, err)
					alive = false
					break
				}
				m.pongs.Add(1)
				lastIO = time.Now()
				if !idleTimer.Stop() {
					select { case <-idleTimer.C: default: }
				}
				idleTimer.Reset(cfg.idleTimeout)
			case <-idleTimer.C:
				// Optional: close idle connection to refresh
				_ = c.Close()
				alive = false
			default:
				_ = lastIO // spin lightly; could block on reads if needed
				time.Sleep(50 * time.Millisecond)
			}
		}
		_ = c.Close()
		m.active.Add(^int64(0))
		m.restarts.Add(1)

		// small backoff before next reconnect
		time.Sleep(jitter(200*time.Millisecond, 0.5))
	}
}

func parseAddrs(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		// add default port if missing
		if !strings.Contains(p, ":") {
			p = p + ":6379"
		}
		out = append(out, p)
	}
	return out
}

func main() {
	var (
		addrList   string
		numConns   int
		username   string
		password   string
		db         int
		useTLS     bool
		pingEvery  time.Duration
		readTO     time.Duration
		idleTO     time.Duration
		connectTO  time.Duration
		openRate   int
	)

	flag.StringVar(&addrList, "addrs", getenv("REDIS_ADDRESSES", "localhost:6379"), "Comma-separated host:port for cluster nodes")
	flag.IntVar(&numConns, "connections", getenvInt("CONNECTIONS", 1000), "Number of concurrent connections")
	flag.StringVar(&username, "username", getenv("REDIS_USERNAME", ""), "ACL username (optional)")
	flag.StringVar(&password, "password", getenv("REDIS_PASSWORD", ""), "Password (optional)")
	flag.IntVar(&db, "db", getenvInt("REDIS_DB", 0), "Database index (optional)")
	flag.BoolVar(&useTLS, "tls", getenvBool("REDIS_TLS", false), "Use TLS for connections")
	flag.DurationVar(&pingEvery, "ping-interval", getenvDuration("PING_INTERVAL", 30*time.Second), "Interval between PINGs per connection")
	flag.DurationVar(&readTO, "read-timeout", getenvDuration("READ_TIMEOUT", 2*time.Second), "Read deadline for responses")
	flag.DurationVar(&idleTO, "idle-timeout", getenvDuration("IDLE_TIMEOUT", 10*time.Minute), "Idle time before recycling a connection")
	flag.DurationVar(&connectTO, "connect-timeout", getenvDuration("CONNECT_TIMEOUT", 5*time.Second), "Dial timeout")
	flag.IntVar(&openRate, "open-rate", getenvInt("OPEN_RATE", 200), "Max new connections per second")
	flag.Parse()

	addrs := parseAddrs(addrList)
	if len(addrs) == 0 {
		log.Fatal("no redis addresses provided")
	}

	targets := make([]dialTarget, 0, len(addrs))
	for _, a := range addrs {
		targets = append(targets, dialTarget{addr: a, useTLS: useTLS})
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	m := &metrics{}

	// Rate-limit opening bursts to avoid SYN floods
	var wg sync.WaitGroup
	tokens := make(chan struct{}, openRate)
	go func() {
		ticker := time.NewTicker(time.Second)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				for i := 0; i < openRate; i++ {
					select {
					case tokens <- struct{}{}:
					default:
					}
				}
			}
		}
	}()

	for i := 0; i < numConns; i++ {
		<-tokens // throttle
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			cfg := connWorkerCfg{
				id:          id,
				targets:     targets,
				username:    username,
				password:    password,
				db:          db,
				pingEvery:   pingEvery,
				readTimeout: readTO,
				idleTimeout: idleTO,
				connectTO:   connectTO,
				backoffMin:  200 * time.Millisecond,
				backoffMax:  5 * time.Second,
			}
			connWorker(ctx, cfg, m)
		}(i)
	}

	// Metrics logger
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		start := time.Now()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				uptime := time.Since(start).Truncate(time.Second)
				log.Printf("uptime=%s active=%d opened=%d authed=%d pongs=%d errors=%d restarts=%d", uptime, m.active.Load(), m.opened.Load(), m.authed.Load(), m.pongs.Load(), m.errors.Load(), m.restarts.Load())
			}
		}
	}()

	wg.Wait()
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getenvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func getenvBool(key string, def bool) bool {
	if v := os.Getenv(key); v != "" {
		switch strings.ToLower(v) {
		case "1", "true", "yes", "y":
			return true
		case "0", "false", "no", "n":
			return false
		}
	}
	return def
}

func getenvDuration(key string, def time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return def
}
