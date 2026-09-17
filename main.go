package main

import (
	"context"
	"fmt"
	"github.com/martin0b101/watcher/config"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type FileSnapshot struct {
	Name       string
	ModifiedAt time.Time
	Size       int64
}

// because the method doesn't modify f, so there's no reason to use a pointer receiver.
func (f FileSnapshot) IsModified(other FileSnapshot) bool {
	return f.ModifiedAt.After(other.ModifiedAt) || f.Size != other.Size
}

func main() {
	config, err := config.Parse(os.Args)
	if err != nil {
		fmt.Printf("error while parsing config: %v\n", err)
	}

	fmt.Printf("\nWatching %s...\n\n", config.Path)

	ctx, stopFunc := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stopFunc()

	watcher, err := createWatcher(ctx, *config)

	if err != nil {
		fmt.Printf("error while running watcher: %v\n", err)
		return
	}
	watcher.Run(ctx)
}

func createWatcher(ctx context.Context, cfg config.Config) (*Watcher, error) {
	watcherTicker, err := NewWatcherTicker(cfg.Interval)
	if err != nil {
		return nil, err
	}

	var reporter Reporter
	switch cfg.OutputFormat {
	case config.Console:
		reporter = ConsoleReporter{}
	case config.JSON:
		fmt.Println("JSON output is not supported right now, fallback to console output")
		reporter = ConsoleReporter{}
	}

	watcher, err := NewWatcher(ctx, cfg.Path, reporter, OSFileSystem{
		Recursive: cfg.Recursive,
	}, watcherTicker)
	if err != nil {
		return nil, err
	}
	return watcher, nil
}
