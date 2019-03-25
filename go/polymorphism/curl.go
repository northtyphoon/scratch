package polymorphism

import (
	"net/http"
	"fmt"
	"strings"
	"io/ioutil"
)

type Curl interface {
	Test(completed chan error)
}

type getCurl struct {
	url string
}

func (g *getCurl) Test(completed chan error) {
	client := &http.Client{}

	req, err:= http.NewRequest("GET", g.url, nil)
	if err != nil {
		completed <- err
	}

	resp, err := client.Do(req)
	if err != nil {
		completed <- err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		completed <- err
	}

	fmt.Printf("Test (Get): %s, Resp: %s", g.url, string(body))
	fmt.Println()

	completed <- nil
}

type postCurl struct {
	url string
	data string
}

func (p *postCurl) Test(completed chan error) {
	client := &http.Client{}

	req, err:= http.NewRequest("POST", p.url, strings.NewReader(p.data))
	if err != nil {
		completed <- err
	}

	resp, err := client.Do(req)
	if err != nil {
		completed <- err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		completed <- err
	}
	
	fmt.Printf("Test (Post): %s with %s, Resp: %s", p.url, p.data, string(body))

	completed <-  nil
}