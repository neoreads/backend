package util

import (
	"log"
	"os"
	"path/filepath"
)

func StripExt(path string) string {
	return path[0 : len(path)-len(filepath.Ext(path))]
}

func FindFile(dir string, pattern string) []string {
	var files []string
	filepath.Walk(dir, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			matched, err := filepath.Match(pattern, filepath.Base(path))
			if err != nil {
				log.Printf("Error in matching pattern %v for file %v", pattern, path)
				return err
			}
			if matched {
				files = append(files, path)
			}
		}
		return nil
	})
	return files
}
