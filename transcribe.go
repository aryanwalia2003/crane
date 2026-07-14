package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type whisperOffsets struct {
	From int64 `json:"from"`
}

type whisperSegment struct {
	Offsets whisperOffsets `json:"offsets"`
	Text    string         `json:"text"`
}

type whisperOutput struct {
	Transcription []whisperSegment `json:"transcription"`
}

func transcribe(whisperBin, model, audio string, speed float64, dir, outPath string) error {
	prefix := filepath.Join(dir, "transcript")
	if err := run(whisperBin, "-m", model, "-f", audio, "-oj", "-of", prefix); err != nil {
		return err
	}
	data, err := os.ReadFile(prefix + ".json")
	if err != nil {
		return err
	}
	var out whisperOutput
	if err := json.Unmarshal(data, &out); err != nil {
		return err
	}
	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, seg := range out.Transcription {
		ts := float64(seg.Offsets.From) / 1000 * speed
		fmt.Fprintf(f, "[%s]%s\n", formatTimestamp(ts), seg.Text)
	}
	return nil
}

func formatTimestamp(sec float64) string {
	h := int(sec) / 3600
	m := (int(sec) % 3600) / 60
	s := int(sec) % 60
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}
