package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
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

	watcherTicker, err := NewWatcherTicker(3*time.Second)
	if err != nil{
		fmt.Printf("error while init ticker %v\n", err)
		return
	}
	watcher, err := NewWatcher(folderPath, ConsoleReporter{}, OSFileSystem{}, watcherTicker)

	if err != nil{
		fmt.Printf("error while running watcher: %v\n", err)
		return
	}

	ctx, stopFunc := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stopFunc()
	watcher.Run(ctx)
}
