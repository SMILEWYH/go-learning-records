//go:build localcontent

// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package website exports the website content for local development.
package website

import (
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
)

// Content returns the website files without embedding the entire site at build
// time. run.sh sets an absolute directory so paths with spaces also work.
func Content() fs.FS {
	dir := os.Getenv("GO_WEBSITE_CONTENT_DIR")
	if dir == "" {
		// Tests run from package directories rather than the repository root.
		_, source, _, ok := runtime.Caller(0)
		if !ok {
			panic("website: set GO_WEBSITE_CONTENT_DIR for localcontent builds")
		}
		dir = filepath.Join(filepath.Dir(source), "_content")
	}
	return os.DirFS(dir)
}

// TourOnly uses the same root as Content: Tour also needs shared images and JS.
func TourOnly() fs.FS {
	return Content()
}
