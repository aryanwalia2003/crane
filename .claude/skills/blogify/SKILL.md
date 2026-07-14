---
name: blogify
description: Turn a crane output directory (transcript.txt + sheet_*.jpg contact sheets + manifest.json) into a complete Hinglish blog post with embedded images pulled from the video. Use when the user asks to turn a crane transcript/manifest into a blog, write up a video as a blog post, or invokes /blogify.
---

Crane (this repo) turns a video into `transcript.txt` (timestamped) plus a set of
`sheet_*.jpg` contact sheets and a `manifest.json` that maps every sampled frame
to `{index, sheet, row, col, timestamp_sec}`. This skill's job is to go from
those raw materials to a finished, complete blog post in Hinglish — picking the
frames worth showing, pulling them out with `crane extract`, and writing prose
around them.

## Inputs

Ask for (or accept as an argument) the crane output directory — default to
`crane-out` in the current directory if one exists. It must contain
`transcript.txt`, `manifest.json`, and one or more `sheet_*.jpg` files. The
`crane` binary must be built (`go build -o crane .` from the repo root) or
already on PATH.

## Steps

1. **Read `transcript.txt` in full.** This is the source of truth for what the
   video actually says — every topic in it must show up somewhere in the blog.
   Do not skip lines because the file is long.

2. **Read `manifest.json`** to get `interval_sec`, `cols`, `rows`, `tile_width`,
   `tile_height`, the sheet list, and the full `frames` array (each entry has
   `index`, `sheet`, `row`, `col`, `timestamp_sec`).

3. **View every `sheet_*.jpg`** with the Read tool (they're plain JPEGs, viewable
   directly). Each sheet is a grid of `cols` x `rows` consecutive frames. Scan
   for frames worth keeping:
   - diagrams, architecture drawings, slides with text, whiteboards, charts, code on screen, product UI
   - skip frames that are just a talking head with no supporting visual
   - skip sponsor/ad segments (recognizable by tone shift in the transcript around the same timestamp) unless the user asks to keep them
   Don't force a quota — a talking-only stretch of the video needs zero images;
   a slide-dense stretch may need several close together.

4. **Map each chosen frame to its global index.** A tile at row `r`, col `c` on
   sheet `s` has index `s*cols*rows + r*cols + c` — or just cross-reference the
   `frames` array in the manifest directly by `sheet`/`row`/`col`.

5. **Extract each chosen frame:**
   ```
   ./crane extract -manifest <dir>/manifest.json -index <N> -out <dir>/blog-images/<slug>.png
   ```
   Use a short descriptive slug per image (not the raw index) so the blog's
   markdown stays readable.

6. **Cross-reference each image's `timestamp_sec`** against the transcript's
   `[HH:MM:SS]` lines to know which part of the narrative it illustrates.

7. **Write `blog.md`** in the crane output directory:
   - Title + short Hinglish intro paragraph.
   - Break the video into sections by topic (not raw timestamp chunks) — use
     headings that match how the video actually flows.
   - Prose in Hinglish for every section. Don't just paraphrase the transcript
     line by line — rewrite it as real blog writing, and add explanation
     wherever a term or mechanism the video mentions quickly deserves more
     unpacking for a reader who wasn't already an expert.
   - Embed each extracted image at the point in the narrative it belongs, via
     `![](blog-images/<slug>.png)`.
   - Completeness check before finishing: walk back through `transcript.txt`
     top to bottom and confirm every topic it raises is represented somewhere
     in the blog. Nothing gets dropped just because it was a short mention.
   - Close with a short wrap-up paragraph.

## Output

`blog.md` and `blog-images/*.png` inside the given crane output directory,
ready to publish as-is.
