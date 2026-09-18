package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/martin0b101/watcher/config"

	"github.com/martin0b101/watcher/internal/filesystem"
	"github.com/martin0b101/watcher/internal/reporter"
	"github.com/martin0b101/watcher/internal/watcher"
)

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

func createWatcher(ctx context.Context, cfg config.Config) (*watcher.Watcher, error) {
	watcherTicker, err := watcher.NewWatcherTicker(cfg.Interval)
	if err != nil {
		return nil, err
	}

	var rpt watcher.Reporter
	switch cfg.OutputFormat {
	case config.Console:
		rpt = reporter.ConsoleReporter{}
	case config.JSON:
		rpt = reporter.ConsoleReporter{}
		rpt.ReportMessage("JSON output is not supported right now, fallback to console output")
	}

	watcher, err := watcher.NewWatcher(ctx, cfg.Path, rpt, filesystem.OSFileSystem{
		Recursive: cfg.Recursive,
	}, watcherTicker)
	if err != nil {
		return nil, err
	}
	return watcher, nil
}
