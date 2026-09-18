package watcher

import (
	"context"
	"errors"
	"fmt"
	"github.com/martin0b101/watcher/internal/filesystem"
	"github.com/martin0b101/watcher/internal/model"
)

type folderState int

const (
	folderStateAvailable folderState = iota
	folderStateUnavailable
)

type Watcher struct {
	path             string
	folderState      folderState
	previousSnapshot map[string]model.FileSnapshot
	reporter         Reporter
	fileSystem       FileSystem
	ticker           Ticker
}

func NewWatcher(ctx context.Context, path string, reporter Reporter, fileSystem FileSystem, ticker Ticker) (*Watcher, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if path == "" {
		return nil, fmt.Errorf("path is empty")
	}

	if reporter == nil {
		return nil, fmt.Errorf("reporter is not init")
	}

	if fileSystem == nil {
		return nil, fmt.Errorf("fileSystem is not init")
	}

	if ticker == nil {
		return nil, fmt.Errorf("ticker is not init")
	}

	snapshot, err := fileSystem.Scan(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize snapshot: %w", err)
	}

	return &Watcher{
		path:             path,
		folderState:      folderStateAvailable,
		previousSnapshot: snapshot,
		reporter:         reporter,
		fileSystem:       fileSystem,
		ticker:           ticker,
	}, nil
}

func (w *Watcher) Run(ctx context.Context) {
	defer w.ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-w.ticker.C():
			currentSnapshot, err := w.fileSystem.Scan(ctx, w.path)
			if err != nil {
				if errors.Is(err, filesystem.ErrFolderUnavailable) && w.folderState != folderStateUnavailable {
					w.folderState = folderStateUnavailable
					w.reporter.ReportMessage(fmt.Sprintf("Directory unavailable: %s", w.path))
				}
				continue
			}
			if w.folderState == folderStateUnavailable {
				w.folderState = folderStateAvailable
				w.reporter.ReportMessage(fmt.Sprintf("Directory available: %s", w.path))
			}
			changes := w.check(currentSnapshot)
			w.reporter.Report(changes)
		}
	}
}

func (w *Watcher) check(currentSnapshot map[string]model.FileSnapshot) (changes model.Changes) {
	changes.Added, changes.Removed, changes.Modified = compare(currentSnapshot, w.previousSnapshot)
	w.previousSnapshot = currentSnapshot
	return changes
}

func compare(currentSnapshot, previousSnapshot map[string]model.FileSnapshot) (added []string, removed []string, modified []string) {
	for previousFile, previousFileContent := range previousSnapshot {
		currentFile, fileExist := currentSnapshot[previousFile]

		if !fileExist {
			removed = append(removed, previousFile)
		} else if currentFile.IsModified(previousFileContent) {
			modified = append(modified, previousFile)
		}
	}
	for file := range currentSnapshot {
		_, fileExist := previousSnapshot[file]
		if !fileExist {
			added = append(added, file)
		}
	}
	return added, removed, modified
}
