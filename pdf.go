package main

import (
	_ "embed"
	"os"
	"path/filepath"
)

//go:embed mdpdf.py
var mdpdfScript string

func convertPDF(in, out string) error {
	tmp, err := os.CreateTemp("", "crane-mdpdf-*.py")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.WriteString(mdpdfScript); err != nil {
		return err
	}
	tmp.Close()
	inAbs, err := filepath.Abs(in)
	if err != nil {
		return err
	}
	outAbs, err := filepath.Abs(out)
	if err != nil {
		return err
	}
	return run("python3", tmp.Name(), inAbs, outAbs)
}
