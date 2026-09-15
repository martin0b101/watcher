package main

import (
	"fmt"
	"io/fs"
	"os"
	"time"
)

func main()  {
	args := os.Args
	if len(args) < 2 {
		fmt.Printf("Folder path is missing.\n")
		return
	}

	folderPath := os.Args[1]

	filesMemory := make(map[string]struct{}, 10)

	err := addInitalFilesToMemory(folderPath, filesMemory)
	if err != nil{
		fmt.Printf("error scaning folder: %v", err)
		return
	}

	fmt.Printf("\nWatching %s...\n\n", folderPath)

	for {
		time.Sleep(3 * time.Second)
		err := scanFolder(folderPath, filesMemory)
		if err != nil{
			fmt.Printf("error scaning folder: %v", err)
		}
	}
}

func scanFolder(folderPath string, memory map[string]struct{}) error{
	files, err := os.ReadDir(folderPath)
	if err != nil{
		return err
	}
	addedFolders := add(files, memory)
	for _, fileName := range addedFolders {		
		fmt.Printf("New file: %s\n", fileName)
	}
	return nil
}

func addInitalFilesToMemory(folderPath string, memory map[string]struct{}) error{
	filenameInFolder, err := os.ReadDir(folderPath)
	if err != nil{
		return err
	}
	add(filenameInFolder, memory)
	return nil
}

func add(files []fs.DirEntry, memory map[string]struct{}) []string{
	var newFiles []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileName := file.Name()
		_, exist := memory[fileName]
		if exist{
			continue
		}
		memory[fileName] = struct{}{}
		newFiles = append(newFiles, fileName)
	}
	return newFiles
}