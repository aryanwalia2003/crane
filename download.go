package main

import (
	"path/filepath"
	"strings"
)

func isURL(s string) bool {
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}

func downloadVideo(ytdlpBin, url, dir string) (string, error) {
	out := filepath.Join(dir, "source.%(ext)s")
	if err := run(ytdlpBin, "-f", "bv*+ba/b", "--merge-output-format", "mp4", "-o", out, url); err != nil {
		return "", err
	}
	return filepath.Join(dir, "source.mp4"), nil
}

func resolveVideo(ytdlpBin, input, dir string) (string, error) {
	if isURL(input) {
		return downloadVideo(ytdlpBin, input, dir)
	}
	return input, nil
}
