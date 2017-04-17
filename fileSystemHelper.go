package main

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"unsafe"
)

func directoryExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

func directoryExistsAll(root string, dirs []staticDir) (invalidPath string, state bool) {
	for _, s := range dirs {
		current := root + s.path
		if !directoryExists(current) {
			return current, false
		}
	}
	return "", true
}

func readFile(p string) (string, error) {
	content, err := ioutil.ReadFile(p)
	if err != nil {
		fmt.Printf("unexcpected error during file read: %s\n", err)
		return "", err
	}

	return string(content), nil
}

func writeFile(fn string, file string) error {
	return ioutil.WriteFile(fn, []byte(file), 0777)
}

func getFileHash(p string) (string, error) {
	content, err := readFile(p)
	if err != nil {
		return "", nil
	}

	sha := sha1.New()
	sha.Write([]byte(content))
	hash := fmt.Sprintf("%x", sha.Sum(nil))

	return hash, nil
}

func strCopy(s string) string {
	var b []byte
	h := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	h.Data = (*reflect.StringHeader)(unsafe.Pointer(&s)).Data
	h.Len = len(s)
	h.Cap = len(s)
	return string(b)
}
