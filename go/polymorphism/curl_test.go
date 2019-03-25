package polymorphism

import (
	"testing"
	"rsc.io/quote"
)

func TestCurl(t *testing.T) {
	tests := []struct{
		curl Curl
		completed chan error
	}{
		{
			&getCurl{
				url: "https://httpbin.org/get",
			},
			make(chan error),
		},
		{
			&postCurl{
				url: "https://httpbin.org/post",
				data: quote.Glass(),
			},
			make(chan error),
		},
	}

	for _, test := range tests {
		go test.curl.Test(test.completed)
	}

	for _, test := range tests {
		err := <- test.completed
		if err != nil {
			t.Error(err)
		}
	}
}