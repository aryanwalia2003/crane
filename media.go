package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

func atempoChain(speed float64) string {
	var parts []string
	for speed > 2.0 {
		parts = append(parts, "atempo=2.0")
		speed /= 2.0
	}
	for speed < 0.5 {
		parts = append(parts, "atempo=0.5")
		speed /= 0.5
	}
	parts = append(parts, fmt.Sprintf("atempo=%.4f", speed))
	return strings.Join(parts, ",")
}

func extractAudio(ffmpegBin, video string, speed float64, dir string) (string, error) {
	out := filepath.Join(dir, "audio.wav")
	err := run(ffmpegBin, "-y", "-i", video, "-vn", "-ac", "1", "-ar", "16000",
		"-filter:a", atempoChain(speed), out)
	return out, err
}

func extractFrames(ffmpegBin, video string, interval float64, tileWidth int, dir string) error {
	pattern := filepath.Join(dir, "frame_%05d.jpg")
	vf := fmt.Sprintf("fps=1/%f,scale=%d:-2", interval, tileWidth)
	return run(ffmpegBin, "-y", "-i", video, "-vf", vf, "-vsync", "0", pattern)
}
