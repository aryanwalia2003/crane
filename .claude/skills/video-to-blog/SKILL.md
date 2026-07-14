---
name: video-to-blog
description: One-shot pipeline from a raw video URL to a finished Hinglish blog PDF, using crane. Use when the user gives a video link and asks for a blog, a blog post, a writeup, or a PDF out of it, or invokes /video-to-blog.
---

This chains crane's whole pipeline end to end: download at 720p, transcribe +
sample frames, write the blog, render the PDF. It's this repo's `blogify`
skill plus the mechanical steps around it.

## Inputs

- A video URL (or local file path) — required, from the user's message or the
  skill argument.
- An output directory — default to `crane-out` in the current directory if
  the user doesn't name one.
- `CRANE_MODEL` must be set (or pass `-model` explicitly) to a local
  distil-whisper GGML model path. If neither is available, stop and ask for
  the model path before doing anything else.

## Steps

1. Make sure `crane` is built (`go build -o crane .` from the repo root if the
   binary is missing or stale).

2. Run:
   ```
   ./crane process -model "$CRANE_MODEL" -out <dir> <video-url-or-path>
   ```
   This downloads at max 720p by default (`-max-height 720`), transcribes with
   local whisper, and produces `<dir>/transcript.txt`, `<dir>/sheet_*.jpg`, and
   `<dir>/manifest.json`. Pass `-max-height` through only if the user asked for
   a different cap.

3. Follow this repo's `blogify` skill (`.claude/skills/blogify/SKILL.md`)
   against `<dir>` to produce `<dir>/blog.md` and `<dir>/blog-images/*.png` —
   read the transcript in full, view every contact sheet, pick the frames
   worth keeping, extract them, and write the complete Hinglish blog post.
   Do not skip this step or shortcut it — it's the part that actually needs
   judgment, and completeness (nothing from the video missing) matters most
   here.

4. Render the PDF:
   ```
   ./crane pdf -in <dir>/blog.md -out <dir>/blog.pdf
   ```

5. Report the final path to `<dir>/blog.pdf` to the user.

## Output

`<dir>/transcript.txt`, `<dir>/sheet_*.jpg`, `<dir>/manifest.json`,
`<dir>/blog.md`, `<dir>/blog-images/*.png`, `<dir>/blog.pdf`.
