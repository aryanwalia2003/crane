package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

func isURL(s string) bool {
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}

func downloadVideo(ytdlpBin, url, dir string, maxHeight int) (string, error) {
	out := filepath.Join(dir, "source.%(ext)s")
	format := fmt.Sprintf("bv*[height<=%d]+ba/b[height<=%d]", maxHeight, maxHeight)
	if err := run(ytdlpBin, "-f", format, "--merge-output-format", "mp4", "-o", out, url); err != nil {
		return "", err
	}
	return filepath.Join(dir, "source.mp4"), nil
}

func resolveVideo(ytdlpBin, input, dir string, maxHeight int) (string, error) {
	if isURL(input) {
		return downloadVideo(ytdlpBin, input, dir, maxHeight)
	}
	return input, nil
}
