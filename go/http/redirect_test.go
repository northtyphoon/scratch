package http

import (
	"fmt"
	"testing"
)

func TestDownloadContent(t *testing.T) {
	registry := "myregistry.azurecr.io"
	repo := "myrepo"
	digest := "sha256:abc"
	uri := "https://" + registry + "/v2/" + repo + "/blobs/" + digest
	authHeader := "Bearer abc"
	content, err := DownloadContent(uri, authHeader)
	if err == nil {
		fmt.Println("succeeded")
		fmt.Println(len(content))
	} else {
		fmt.Println("failed")
	}
}
