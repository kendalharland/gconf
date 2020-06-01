package github

type GetTreeResponse struct {
	SHA       string     `json:"sha"`
	URL       string     `json:"url"`
	Tree      []TreeNode `json:"tree"`
	Truncated bool       `json:"truncated"`
}

type TreeNode struct {
	Path string `json:"path"`
	Mode string `json:"mode"`
	Size int64  `json:"size"`
	SHA  string `json:"sha"`
	Type string `json:"type"`
	URL  string `json:"url"`
}

// https://developer.github.com/v3/git/blobs
type GetBlobResponse struct {
	Content  string `json:"content"`
	Encoding string `json:"encoding"`
	URL      string `json:"url"`
	SHA      string `json:"sha"`
	Size     int64  `json:"size"`
}
