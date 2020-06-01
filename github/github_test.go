package github

import (
	"encoding/json"
	"github.com/google/go-cmp/cmp"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetTree(t *testing.T) {
	message := &GetTreeResponse{
		SHA: "fc6274d15fa3ae2ab983129fb037999f264ba9a7",
		URL: "https://api.github.com/repos/owner/repo/trees/sha",
		Tree: []TreeNode{{
			Path: "subdir/file.txt",
			Mode: "100644",
			Type: "blob",
			Size: 132,
			SHA:  "7c258a9869f33c1e1e1f74fbb32f07c86cb5a75b",
			URL:  "https://api.github.com/repos/owner/repo/trees/sha",
		}},
		Truncated: false,
	}

	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(message); err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()

	res, err := GetTree(ts.Client(), ts.URL, "owner", "repo", "sha", true)
	if err != nil {
		t.Error(err)
	}

	got := res
	want := message
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("GetTree() mismatch (-want +got):\n%s", diff)
	}
}

func TestGetBlob(t *testing.T) {
	message := &GetBlobResponse{
		SHA:      "fc6274d15fa3ae2ab983129fb037999f264ba9a7",
		URL:      "https://api.github.com/repos/owner/repo/trees/sha",
		Encoding: "base64",
		Size:     5,
		Content:  "hello",
	}
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(message); err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()

	res, err := GetBlob(ts.Client(), ts.URL, "owner", "repo", "sha")
	if err != nil {
		t.Error(err)
	}

	got := res
	want := message
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("GetBlob() mismatch (-want +got):\n%s", diff)
	}
}
