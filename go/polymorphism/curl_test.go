package polymorphism

import (
	"testing"
	"rsc.io/quote"
)

func TestCurl(t *testing.T) {
	curls := []Curl{
		&getCurl{
			url: "https://httpbin.org/get",
		},
		&postCurl{
			url: "https://httpbin.org/post",
			data: quote.Glass(),
		},
	}

	for _, curl := range curls {
		if err := curl.Test(); err != nil {
			t.Error(err)
		}
	}
}