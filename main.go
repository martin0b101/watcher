package main

import (
	"fmt"
	"os"
	"time"
)

type FileSnapshot struct {
	Name string
	ModifiedAt time.Time
	Size int64
} 

//because the method doesn't modify f, so there's no reason to use a pointer receiver.
func (f FileSnapshot) IsModified(other FileSnapshot) bool{
	return f.ModifiedAt.After(other.ModifiedAt) || f.Size != other.Size
}

func main()  {
	if len(os.Args) < 2 {
		fmt.Printf("Folder path is missing.\n")
		return
	}
	folderPath := os.Args[1]

	fmt.Printf("\nWatching %s...\n\n", folderPath)

	reporter := NewConsoleReporter(nil)
	watcher, err := NewWatcher(folderPath, 3*time.Second, reporter, OSFileSystem{})

	if err != nil{
		fmt.Printf("error while running watcher: %v/n", err)
		return
	}
	watcher.Run()
}
