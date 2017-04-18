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

var permissions map[string]os.FileMode

func readFile(path string) (fc string, err error) {
	if permissions == nil {
		permissions = make(map[string]os.FileMode)
	}

	info, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("unexpected error during file permission read: %s", err)
		return
	}
	content, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("unexcpected error during file read: %s\n", err)
		return "", err
	}

	permissions[path] = info.Mode()

	return string(content), nil
}

func writeFile(path string, file string) error {
	return ioutil.WriteFile(path, []byte(file), permissions[path])
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
