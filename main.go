package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

func main() {
	debugLev = debugNo

	root, htmlRegex, sIn, err := organizeCLArguments(os.Args...)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if !directoryExists(root) {
		fmt.Printf("root directory (%s) doesn't exist\n", root)
		os.Exit(1)
	}

	path, state := directoryExistsAll(root, sIn)
	if !state {
		fmt.Printf("static asset directory (%s) doesn't exist\n", path)
		os.Exit(1)
	}

	htmlFiles, err := searchHTMLFiles(root, htmlRegex)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	staticFiles, err := searchStaticFiles(root, sIn)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var concurrentStaticFileMap = struct {
		sync.RWMutex
		files map[string]string
	}{files: staticFiles}

	var wg sync.WaitGroup
	for fn, file := range htmlFiles {
		wg.Add(1)
		go func(fn string, file string) {
			defer wg.Done()

			fn = root + strings.TrimPrefix(fn, "/")
			fmt.Printf("-> %s\n", fn)

			ferr := analyzeFile(fn, file, &concurrentStaticFileMap)
			if ferr != nil {
				fmt.Println(ferr)
				os.Exit(1)
			}
		}(fn, file)
	}

	for fn, newFn := range staticFiles {
		err = os.Rename(root+strings.TrimPrefix(fn, "/"), root+strings.TrimPrefix(newFn, "/"))
		if err != nil {
			log.Println(err.Error())
		}
	}

	wg.Wait()
}
