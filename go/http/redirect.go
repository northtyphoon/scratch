package http

import (
	"fmt"
	"io"
	"net/http"
)

func DownloadContent(uri, authHeader string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		fmt.Printf("%v", err)
	}
	req.Header.Add("Authorization", authHeader)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("%v", err)
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
