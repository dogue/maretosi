package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// removes the given root and all ancestors from a path
func extractSubpath(path, root string) (subpath string) {
	root = filepath.Base(root)
	root += "/"
	_, subpath, _ = strings.Cut(path, root)
	return
}

func mdToHtml(path string) string {
	return strings.Replace(path, ".md", ".html", 1)
}

// checks that a path exists and is a directory
//
// returns nil if the checks pass, else an error
func validDir(path string) error {
	path, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to determine absolute path for file %q: %w", path, err)
	}

	fileInfo, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to get file info for %q: %w", path, err.(*fs.PathError).Err)
	}

	if !fileInfo.IsDir() {
		return fmt.Errorf("path is not a directory: %q", path)
	}

	return nil
}

func validateDirs() error {
	if err := validDir(contentDir); err != nil {
		return err
	}

	if err := validDir(outputDir); err != nil {
		err = os.MkdirAll(outputDir, fs.ModePerm)
		if err != nil {
			return err
		}
	}

	if err := validDir(assetsDir); err != nil && !skipAssets {
		return err
	}

	return nil

}
