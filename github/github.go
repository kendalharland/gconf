// Package Github implements a client for v3 of the Github API.
//
// For full documentation of the API see: https://developer.github.com/v3
package github

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const DefaultHost = "https://api.github.com"

func GetTree(client *http.Client, host, owner, repo, treeSHA string, recursive bool) (*GetTreeResponse, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/git/trees/%s?recursive=%t", host, owner, repo, treeSHA, recursive)
	res, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	message := new(GetTreeResponse)
	if err := decodeResponse(res, message); err != nil {
		return nil, err
	}
	return message, nil
}

func GetBlob(client *http.Client, host, owner, repo, fileSHA string) (*GetBlobResponse, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/git/blobs/%s", host, owner, repo, fileSHA)
	res, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	message := new(GetBlobResponse)
	if err := decodeResponse(res, message); err != nil {
		return nil, err
	}
	return message, nil
}

func decodeResponse(r *http.Response, into interface{}) error {
	defer r.Body.Close()

	if r.StatusCode != 200 {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return fmt.Errorf("%s %w", r.Status, err)
		}
		return fmt.Errorf("%s %s", r.Status, string(body))
	}
	if r.Header.Get("X-Ratelimit-Remaining") == "0" {
		maxPerHour := r.Header.Get("X-RateLimit-Limit")
		return fmt.Errorf("Github API quota exceeded %s requests per hour", maxPerHour)
	}

	return json.NewDecoder(r.Body).Decode(into)
}
