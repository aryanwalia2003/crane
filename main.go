package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}
	var err error
	switch os.Args[1] {
	case "process":
		err = cmdProcess(os.Args[2:])
	case "extract":
		err = cmdExtract(os.Args[2:])
	default:
		usage()
		os.Exit(1)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: crane <process|extract> [flags]")
}

func cmdProcess(args []string) error {
	fs := flag.NewFlagSet("process", flag.ExitOnError)
	outDir := fs.String("out", "crane-out", "output directory")
	interval := fs.Float64("interval", 5, "seconds between sampled frames")
	cols := fs.Int("cols", 5, "sprite sheet columns")
	rows := fs.Int("rows", 5, "sprite sheet rows")
	tileWidth := fs.Int("tile-width", 320, "frame tile width in pixels")
	speed := fs.Float64("speed", 2, "audio speedup factor before transcription")
	whisperBin := fs.String("whisper-bin", "whisper-cli", "whisper.cpp binary")
	model := fs.String("model", "", "path to ggml distil-whisper model (required)")
	ytdlpBin := fs.String("yt-dlp-bin", "yt-dlp", "yt-dlp binary")
	ffmpegBin := fs.String("ffmpeg-bin", "ffmpeg", "ffmpeg binary")
	fs.Parse(args)
	if fs.NArg() < 1 {
		return fmt.Errorf("usage: crane process [flags] <url-or-file>")
	}
	if *model == "" {
		return fmt.Errorf("-model is required")
	}
	return runProcess(processOpts{
		input:      fs.Arg(0),
		outDir:     *outDir,
		interval:   *interval,
		cols:       *cols,
		rows:       *rows,
		tileWidth:  *tileWidth,
		speed:      *speed,
		whisperBin: *whisperBin,
		model:      *model,
		ytdlpBin:   *ytdlpBin,
		ffmpegBin:  *ffmpegBin,
	})
}

func cmdExtract(args []string) error {
	fs := flag.NewFlagSet("extract", flag.ExitOnError)
	manifestPath := fs.String("manifest", "manifest.json", "path to manifest.json")
	index := fs.Int("index", -1, "frame index to extract")
	out := fs.String("out", "", "output image path (required)")
	fs.Parse(args)
	if *index < 0 {
		return fmt.Errorf("-index is required")
	}
	if *out == "" {
		return fmt.Errorf("-out is required")
	}
	return extractFrame(*manifestPath, *index, *out)
}
