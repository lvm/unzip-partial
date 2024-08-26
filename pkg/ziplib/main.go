package ziplib

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Zip interface {
	Extract(pattern, outputDir string, fileCh chan<- string, errCh chan<- error)
	List(fileCh chan<- string, errCh chan<- error)
}

type ZipFile struct {
	Path string
}

func (z *ZipFile) Extract(pattern, outputDir string, fileCh chan<- string, errCh chan<- error) {
	r, err := zip.OpenReader(z.Path)
	if err != nil {
		errCh <- fmt.Errorf("failed to open zip file: %w", err)
		return
	}
	defer r.Close()

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		errCh <- fmt.Errorf("failed to create output directory: %w", err)
		return
	}

	for _, f := range r.File {
		matched, err := filepath.Match(pattern, f.Name)
		if err != nil {
			errCh <- fmt.Errorf("failed to match pattern: %w", err)
			return
		}

		if matched {
			rc, err := f.Open()
			if err != nil {
				errCh <- fmt.Errorf("failed to open file in archive: %w", err)
				return
			}
			defer rc.Close()

			outputPath := filepath.Join(outputDir, f.Name)
			if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
				errCh <- fmt.Errorf("failed to create directories for output file: %w", err)
				return
			}

			outFile, err := os.Create(outputPath)
			if err != nil {
				errCh <- fmt.Errorf("failed to create output file: %w", err)
				return
			}
			defer outFile.Close()

			if _, err = io.Copy(outFile, rc); err != nil {
				errCh <- fmt.Errorf("failed to copy file contents: %w", err)
				return
			}

			fileCh <- f.Name
		}
	}
}

func (z *ZipFile) List(fileCh chan<- string, errCh chan<- error) {
	r, err := zip.OpenReader(z.Path)
	if err != nil {
		errCh <- fmt.Errorf("failed to open zip file: %w", err)
		return
	}
	defer r.Close()

	for _, f := range r.File {
		fileCh <- f.Name
	}
}
