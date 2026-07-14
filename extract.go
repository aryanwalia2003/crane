package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
)

func extractFrame(manifestPath string, index int, outPath string) error {
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return err
	}
	var m manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	var entry *frameEntry
	for i := range m.Frames {
		if m.Frames[i].Index == index {
			entry = &m.Frames[i]
			break
		}
	}
	if entry == nil {
		return fmt.Errorf("no frame with index %d", index)
	}
	sheetPath := filepath.Join(filepath.Dir(manifestPath), m.Sheets[entry.Sheet])
	img, err := loadImage(sheetPath)
	if err != nil {
		return err
	}
	x0, y0 := entry.Col*m.TileWidth, entry.Row*m.TileHeight
	src := image.Rect(x0, y0, x0+m.TileWidth, y0+m.TileHeight)
	cropped := image.NewRGBA(image.Rect(0, 0, m.TileWidth, m.TileHeight))
	draw.Draw(cropped, cropped.Bounds(), img, src.Min, draw.Src)
	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, cropped)
}
