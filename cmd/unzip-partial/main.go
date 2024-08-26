package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/lvm/unzip-partial/pkg/ziplib"
)

func main() {
	path := flag.String("zip", "", "Path to the zip file")
	pattern := flag.String("pattern", "", "Glob pattern to match files (e.g., '*.doc')")
	output := flag.String("output", "", "Directory to extract matched files to")
	flag.Parse()

	if *path == "" || *pattern == "" || *output == "" {
		fmt.Println("Usage: unzip-partial -zip <zip-file> -pattern <pattern> -output <output-dir>")
		os.Exit(1)
	}

	file := &ziplib.ZipFile{Path: *path}
	extractedFiles := make(chan string)
	errCh := make(chan error)

	go func() {
		defer close(extractedFiles)
		defer close(errCh)
		file.Extract(*pattern, *output, extractedFiles, errCh)
	}()

	for {
		select {
		case fileName, ok := <-extractedFiles:
			if ok {
				log.Printf("Extracted file: %s\n", fileName)
			} else {
				extractedFiles = nil
			}
		case err, ok := <-errCh:
			if ok {
				log.Fatalf("Error extracting files: %v", err)
			} else {
				errCh = nil
			}
		}

		if extractedFiles == nil && errCh == nil {
			break
		}
	}
}
