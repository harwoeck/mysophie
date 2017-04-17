package main

import (
	"fmt"
	"regexp"
	"strings"
)

func organizeCLArguments(args ...string) (root string, htmlRegex *regexp.Regexp, sIn []staticDir, err error) {

	for i := 1; i < len(args); i++ {
		if shouldDebug(debugNormal) {
			fmt.Printf("os.Args[%d] = '%s'\n", i, args[i])
		}

		switch args[i] {
		// handle root
		case "--root":
			fallthrough
		case "-r":
			if len(args) <= i+1 {
				err = fmt.Errorf("invalid usage of '%s' argument. requires one option: 'root-directory-of-web-folder'. %s", args[i], errorMySofieHelp)
				return
			}
			if root != "" {
				err = fmt.Errorf("argument '%s' allowed only once. %s", args[i], errorMySofieHelp)
				return
			}
			root = strings.Replace(args[i+1], "\\", "/", -1)
			if !strings.HasSuffix(root, "/") {
				root += "/"
			}
			i++
			break
		// regex pattern used for matching file names in web root directory
		case "--html":
			fallthrough
		case "-h":
			if len(args) <= i+1 {
				err = fmt.Errorf("invalid usage of '%s' argument. requires one option: 'regex-file-matching-pattern'. %s", args[i], errorMySofieHelp)
				return
			}
			if htmlRegex != nil {
				err = fmt.Errorf("argument '%s' allowed only once. %s", args[i], errorMySofieHelp)
				return
			}
			htmlRegex = regexp.MustCompile(args[i+1])
			i++
			break
		// directory containing static assets files (match with regex option) that
		case "--static":
			fallthrough
		case "-s":
			if len(args) <= i+2 {
				err = fmt.Errorf("invalid usage of '%s' argument. requires two options: 'path-to-assets-directory' & 'regex-file-matching-pattern'. %s", args[i], errorMySofieHelp)
				return
			}
			sIn = append(sIn, *&staticDir{
				path:   args[i+1],
				regstr: args[i+2],
				regex:  regexp.MustCompile(args[i+2]),
			})
			i += 2
			break
		// print detailed help menu with examples
		case "--help":
			fmt.Println("in --help")
			break
		// argument not handled till now -> hence invalid
		default:
			err = fmt.Errorf("invalid argument: '%s'. %s", args[i], errorMySofieHelp)
			return
		}
	}

	// Return error if no root is provided
	if root == "" {
		err = fmt.Errorf("argument '--root' not specified but required once. %s", errorMySofieHelp)
		return
	}

	// If no html file-matching-pattern is set we use the default value (matching a files with .html extension)
	if htmlRegex == nil {
		htmlRegex = regexp.MustCompile(".+\\.html")
	}

	if sIn == nil {
		err = fmt.Errorf("no '--static' argument specified hence mysophie is useless. %s", errorMySofieHelp)
		return
	}

	return
}
