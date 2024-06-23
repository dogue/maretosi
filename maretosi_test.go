package main

import (
	"testing"
)

func TestExtractSubpath(t *testing.T) {
	path := "/home/user/project/content/blog/post.md"
	root := "content"
	subpath := extractSubpath(path, root)

	if subpath != "blog/post.md" {
		t.Fatalf("wrong subpath returned: expected=%q, got=%q", "blog/post.md", subpath)
	}
}
