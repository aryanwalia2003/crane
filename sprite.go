package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"os"
	"path/filepath"
	"sort"
)

type frameEntry struct {
	Index        int     `json:"index"`
	Sheet        int     `json:"sheet"`
	Row          int     `json:"row"`
	Col          int     `json:"col"`
	TimestampSec float64 `json:"timestamp_sec"`
}

type manifest struct {
	IntervalSec float64      `json:"interval_sec"`
	TileWidth   int          `json:"tile_width"`
	TileHeight  int          `json:"tile_height"`
	Cols        int          `json:"cols"`
	Rows        int          `json:"rows"`
	Sheets      []string     `json:"sheets"`
	Frames      []frameEntry `json:"frames"`
}

func buildSprite(framesDir, outDir string, interval float64, cols, rows int) (*manifest, error) {
	files, err := filepath.Glob(filepath.Join(framesDir, "frame_*.jpg"))
	if err != nil {
		return nil, err
	}
	sort.Strings(files)
	if len(files) == 0 {
		return nil, fmt.Errorf("no frames extracted")
	}
	first, err := loadImage(files[0])
	if err != nil {
		return nil, err
	}
	tw, th := first.Bounds().Dx(), first.Bounds().Dy()
	perSheet := cols * rows
	m := &manifest{IntervalSec: interval, TileWidth: tw, TileHeight: th, Cols: cols, Rows: rows}
	var canvas *image.RGBA
	sheetIdx := -1
	for i, path := range files {
		pos := i % perSheet
		if pos == 0 {
			if canvas != nil {
				if err := saveSheet(canvas, outDir, sheetIdx); err != nil {
					return nil, err
				}
			}
			sheetIdx++
			canvas = image.NewRGBA(image.Rect(0, 0, tw*cols, th*rows))
			m.Sheets = append(m.Sheets, sheetName(sheetIdx))
		}
		img, err := loadImage(path)
		if err != nil {
			return nil, err
		}
		row, col := pos/cols, pos%cols
		dst := image.Rect(col*tw, row*th, col*tw+tw, row*th+th)
		draw.Draw(canvas, dst, img, image.Point{}, draw.Src)
		m.Frames = append(m.Frames, frameEntry{
			Index: i, Sheet: sheetIdx, Row: row, Col: col,
			TimestampSec: float64(i) * interval,
		})
	}
	if err := saveSheet(canvas, outDir, sheetIdx); err != nil {
		return nil, err
	}
	return m, nil
}

func sheetName(i int) string {
	return fmt.Sprintf("sheet_%03d.jpg", i)
}

func saveSheet(canvas *image.RGBA, outDir string, idx int) error {
	f, err := os.Create(filepath.Join(outDir, sheetName(idx)))
	if err != nil {
		return err
	}
	defer f.Close()
	return jpeg.Encode(f, canvas, &jpeg.Options{Quality: 90})
}

func loadImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	return img, err
}

func writeManifest(m *manifest, path string) error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
