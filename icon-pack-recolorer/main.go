package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
)

var hexRegex = regexp.MustCompile("#([a-f0-9]{6}|[a-f0-9]{3})")

func visit(path string, di fs.DirEntry, err error) error {
	if filepath.Ext(path) != ".svg" {
		return nil
	}
	fileInfo, err := os.Lstat(path)
	if err != nil {
		fmt.Printf("Error getting file info for %s: %v\n", path, err)
	}
	if fileInfo.Mode()&os.ModeSymlink == 0 {
		return nil
	}

	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error getting file contens for %s: %v\n", path, err)
	}

	fileText := string(content)

	editedSVG := hexRegex.ReplaceAllStringFunc(fileText, ClosestMatchInPalette)

	file, err := os.OpenFile(path, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return nil
	}
	defer file.Close()

	_, err = file.WriteString(editedSVG)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		return nil
	}
	return nil
}

func main() {
	flag.Parse()
	root := flag.Arg(0)
	err := filepath.WalkDir(root, visit)
	if err != nil {
		fmt.Println(err)
	}
}
