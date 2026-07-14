# crane

Local video transcript + deduplicatable frame contact sheets, in Go.

Given a video link (or local file), `crane process` runs local whisper
transcription and frame sampling side by side, then tiles the sampled
frames into JPEG contact sheets with a manifest so an LLM can review
them and hand back which frame indices matter (e.g. slides/figures in
a teaching video). `crane extract` then turns a manifest index back
into a standalone image.

## Prerequisites

- `ffmpeg` on PATH
- `yt-dlp` on PATH (only needed for URL input)
- a whisper.cpp CLI binary (`whisper-cli`, formerly `main`) built from
  https://github.com/ggml-org/whisper.cpp, plus a distil-whisper GGML
  model (e.g. `ggml-distil-medium.en.bin`)

## Usage

```
crane process -model ggml-distil-medium.en.bin -out out/ <url-or-file>
```

Outputs into `out/`:
- `transcript.txt` — timestamped transcript
- `sheet_000.jpg`, `sheet_001.jpg`, ... — frame contact sheets
- `manifest.json` — index -> {sheet, row, col, timestamp_sec}

```
crane extract -manifest out/manifest.json -index 7 -out frame_7.png
```

## Flags

`process`: `-interval` (seconds between sampled frames, default 5),
`-cols`/`-rows` (contact sheet grid, default 5x5), `-tile-width`
(default 320px), `-speed` (audio speedup before transcription,
default 2x), `-whisper-bin`, `-model`, `-yt-dlp-bin`, `-ffmpeg-bin`.

`extract`: `-manifest`, `-index`, `-out`.
