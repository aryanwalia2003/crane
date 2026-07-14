package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type processOpts struct {
	input      string
	outDir     string
	interval   float64
	cols, rows int
	tileWidth  int
	speed      float64
	whisperBin string
	model      string
	ytdlpBin   string
	ffmpegBin  string
}

func runProcess(opts processOpts) error {
	if err := os.MkdirAll(opts.outDir, 0755); err != nil {
		return err
	}
	tmp, err := os.MkdirTemp("", "crane-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	video, err := resolveVideo(opts.ytdlpBin, opts.input, tmp)
	if err != nil {
		return fmt.Errorf("resolve video: %w", err)
	}

	var wg sync.WaitGroup
	var transcribeErr, framesErr error

	wg.Add(2)
	go func() {
		defer wg.Done()
		audio, err := extractAudio(opts.ffmpegBin, video, opts.speed, tmp)
		if err != nil {
			transcribeErr = fmt.Errorf("extract audio: %w", err)
			return
		}
		transcribeErr = transcribe(opts.whisperBin, opts.model, audio, opts.speed, tmp,
			filepath.Join(opts.outDir, "transcript.txt"))
	}()
	go func() {
		defer wg.Done()
		framesDir := filepath.Join(tmp, "frames")
		if err := os.MkdirAll(framesDir, 0755); err != nil {
			framesErr = err
			return
		}
		if err := extractFrames(opts.ffmpegBin, video, opts.interval, opts.tileWidth, framesDir); err != nil {
			framesErr = fmt.Errorf("extract frames: %w", err)
			return
		}
		m, err := buildSprite(framesDir, opts.outDir, opts.interval, opts.cols, opts.rows)
		if err != nil {
			framesErr = fmt.Errorf("build sprite: %w", err)
			return
		}
		framesErr = writeManifest(m, filepath.Join(opts.outDir, "manifest.json"))
	}()
	wg.Wait()

	if transcribeErr != nil {
		return transcribeErr
	}
	return framesErr
}
