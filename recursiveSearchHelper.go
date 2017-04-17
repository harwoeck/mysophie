package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func searchHTMLFiles(root string, htmlRegex *regexp.Regexp) (htmlFiles map[string]string, err error) {
	htmlFiles = make(map[string]string)

	// search through complete root directory
	err = filepath.Walk(root, func(p string, f os.FileInfo, err error) error {

		if f.IsDir() {
			return nil
		}

		if htmlRegex.MatchString(f.Name()) {
			file, err := readFile(p)
			if err != nil {
				return err
			}

			htmlFiles[strings.TrimPrefix(p, root)] = file
		}

		return nil
	})

	return
}

func searchStaticFiles(root string, sIn []staticDir) (staticFiles map[string]string, err error) {
	staticFiles = make(map[string]string)

	// Search through all static directories
	for _, dir := range sIn {
		currentSearchDirectory := root + dir.path

		err = filepath.Walk(currentSearchDirectory, func(p string, f os.FileInfo, err error) error {

			if f.IsDir() {
				return nil
			}

			if dir.regex.MatchString(f.Name()) {
				hash, err := getFileHash(p)
				if err != nil {
					return err
				}
				hash = hash[:10]

				p = strings.Replace(p, "\\", "/", -1)
				file := "/" + strings.TrimPrefix(p, root)

				if shouldDebug(debugNormal) {
					fmt.Printf("%s -> '%s'\n", file, hash)
				}

				lastDotIndex := strings.LastIndex(file, ".")
				filename := file[:lastDotIndex]
				extension := file[lastDotIndex:]

				newFile := filename + "-" + hash + extension

				staticFiles[file] = newFile
			}

			return nil
		})

		if err != nil {
			return
		}
	}

	return
}
