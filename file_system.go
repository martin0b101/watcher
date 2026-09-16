package main

import "os"

type FileSystem interface {
	Scan(path string) (map[string]FileSnapshot, error)
}

type OSFileSystem struct {
}

var _ FileSystem = OSFileSystem{}


func (OSFileSystem) Scan(path string) (map[string]FileSnapshot, error) {
	currentSnapshot := make(map[string]FileSnapshot)
	files, err := os.ReadDir(path)
	if err != nil{
		return nil, err
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		info, err := file.Info()
		if err != nil {
			return nil, err
		}
		fileName := file.Name()
		currentSnapshot[fileName] = FileSnapshot{
			Name: fileName,
			ModifiedAt: info.ModTime(),
			Size: info.Size(),
		}
	}
	return currentSnapshot, nil
}
