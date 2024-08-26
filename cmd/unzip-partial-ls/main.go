package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/lvm/unzip-partial/pkg/ziplib"
)

func main() {
	path := flag.String("zip", "", "Path to the zip file to list contents of")
	flag.Parse()

	if *path == "" {
		fmt.Println("Usage: unzip-partial-ls -zip <zip-file>")
		os.Exit(1)
	}

	file := &ziplib.ZipFile{Path: *path}
	listedFiles := make(chan string)
	errCh := make(chan error)

	go func() {
		defer close(listedFiles)
		defer close(errCh)
		file.List(listedFiles, errCh)
	}()

	for {
		select {
		case fileName, ok := <-listedFiles:
			if !ok {
				listedFiles = nil
			} else {
				log.Println(fileName)
			}
		case err, ok := <-errCh:
			if !ok {
				errCh = nil
			} else if err != nil {
				log.Fatalf("Error listing zip contents: %v", err)
			}
		}

		if listedFiles == nil && errCh == nil {
			break
		}
	}
}
