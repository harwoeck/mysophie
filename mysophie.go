package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

func analyzeFile(fn string, file string, staticFiles *struct {
	sync.RWMutex
	files map[string]string
}) error {

	pattern := "data-mysophie=\"\" "
	mustSkip := 1
	changes := 0

	if shouldDebug(debugDetailed) {
		fmt.Printf("initial content of file '%s' is: '%s'\n", fn, file)
	}

	for true {

		fc := strCopy(file)
		startIndex := 0

		// find next occerence of pattern 'data-mysophie="" '
		for i := 0; i < mustSkip; i++ {
			tempIndex := strings.Index(fc, pattern)
			if tempIndex == -1 {
				if shouldDebug(debugNormal) {
					fmt.Println("no more search patterns found - aborting")
				}
				startIndex = tempIndex
				break
			}
			fc = fc[tempIndex+len(pattern):]
			startIndex += tempIndex

			if shouldDebug(debugDetailed) {
				fmt.Printf("[%s] Round(%d): tempIndex=%d | startIndex=%d | fc='%s'\n", fn, i, tempIndex, startIndex, fc)
			}
		}
		if startIndex == -1 {
			break
		}

		endIndex := startIndex
		endIndex += len(pattern)

		// determine what link we have
		isSrc := strings.HasPrefix(fc, "src=\"")
		isHref := strings.HasPrefix(fc, "href=\"")

		// skip if this mysophie pattern was used invalid
		if !isSrc && !isHref {
			fmt.Printf("found mysophie search pattern ('data-mysophie=\"\" ') not followed by src or href: %s:%d\n", fn, endIndex)
			mustSkip++
			continue
		}

		// remove 'script="' or 'href="'
		linkStart := strings.Index(fc, "\"") + 1
		endIndex += linkStart
		fc = fc[linkStart:]

		// remove everything after the next '"'
		linkEnd := strings.Index(fc, "\"")
		endIndex += linkEnd
		fc = fc[:linkEnd]

		if shouldDebug(debugNormal) {
			fmt.Printf("sliced-link-name: '%s' | startIndex: %d | endIndex: %d\n", fc, startIndex, endIndex)
		}

		// build our new link
		staticFiles.RLock()
		newFn := staticFiles.files[fc]
		staticFiles.RUnlock()

		newLink := ""
		if isSrc {
			newLink = "src=\""
		} else {
			newLink = "href=\""
		}
		newLink += newFn + "\""

		if shouldDebug(debugMinimal) {
			fmt.Printf("genarted link: '%s'\n", newLink)
		}

		// find parts before and after our link
		before := file[:startIndex]
		after := file[endIndex+1:]
		if shouldDebug(debugDetailed) {
			fmt.Printf("beforeLink: '%s'\n", before)
			fmt.Printf("afterLink: '%s'\n", after)
		}

		file = before + newLink + after
		changes++

		if shouldDebug(debugDetailed) {
			fmt.Printf("updated file: '%s'\n", file)
		}

		if shouldDebug(debugMinimal) {
			fmt.Println("===\nFinished this file - Press enter for next one")
			reader := bufio.NewReader(os.Stdin)
			reader.ReadString('\n')
		}
	}

	if changes == 0 {
		if shouldDebug(debugMinimal) {
			fmt.Printf("no changes made - file ignored '%s'\n", fn)
		}
		return nil
	}

	if shouldDebug(debugNormal) {
		fmt.Printf("%d changes made - writing file '%s'\n", changes, fn)
	}

	err := writeFile(fn, file)
	if err != nil {
		return err
	}

	return nil
}
