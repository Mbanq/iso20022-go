//go:build ignore
// +build ignore

package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Define the replacements
	replacements := []struct {
		old string
		new string
	}{
		{"innerXml string", "InnerXml string"},
	}

	// Walk through the directory
	err := filepath.WalkDir("ISO20022", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// Only process regular files with .go extension
		if !d.IsDir() && strings.HasSuffix(strings.ToLower(path), ".go") {
			if err := replaceInFile(path, replacements); err != nil {
				fmt.Printf("Error processing %s: %v\n", path, err)
			}
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
	}
}

func replaceInFile(filePath string, replacements []struct{ old, new string }) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	for _, r := range replacements {
		content = bytes.Replace(content, []byte(r.old), []byte(r.new), -1)
	}

	// Write the modified content back to the file
	return os.WriteFile(filePath, content, 0644)
}
