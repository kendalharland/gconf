// Package Github implements a client for v3 of the Github API.
//
// For full documentation of the API see: https://developer.github.com/v3
package github

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const DefaultHost = "https://api.github.com"

func GetTree(client *http.Client, host, owner, repo, treeSHA string, recursive bool) (*GetTreeResponse, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/git/trees/%s?recursive=%t", host, owner, repo, treeSHA, recursive)
	log.Printf("GET %s", url)
	res, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	message := new(GetTreeResponse)
	if err := json.NewDecoder(res.Body).Decode(message); err != nil {
		return nil, err
	}
	return message, nil
}

func GetBlob(client *http.Client, host, owner, repo, fileSHA string) (*GetBlobResponse, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/git/blobs/%s", host, owner, repo, fileSHA)
	log.Printf("GET %s", url)
	res, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	message := new(GetBlobResponse)
	if err := json.NewDecoder(res.Body).Decode(message); err != nil {
		return nil, err
	}
	return message, nil
}
