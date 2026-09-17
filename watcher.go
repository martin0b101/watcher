package main

import (
	"context"
	"errors"
	"fmt"
	"os"
)

type FolderState int 

const (
	FolderStateAvailable FolderState = iota
	FolderStateUnavailable
)


type Watcher struct {
	Path string
	FolderState FolderState
	PreviousSnapshot map[string]FileSnapshot
	Reporter Reporter
	FileSystem FileSystem
	Ticker Ticker
}

func NewWatcher(ctx context.Context,path string, reporter Reporter, fileSystem FileSystem, ticker Ticker) (*Watcher, error) {
	if path == ""{
		return nil, fmt.Errorf("path is empty")
	}

	if reporter == nil{
		return nil, fmt.Errorf("reporter is not init")
	}

	if fileSystem == nil {
		return nil, fmt.Errorf("fileSystem is not init")
	}

	if ticker == nil {
		return nil, fmt.Errorf("ticker is not init")
	}

	snapshot, err := fileSystem.Scan(ctx, path)
	if err != nil{
		return nil, fmt.Errorf("failed to initialize snapshot: %w", err)
	}

	return &Watcher{
		Path: path,
		FolderState: FolderStateAvailable,
		PreviousSnapshot: snapshot,
		Reporter: reporter,
		FileSystem: fileSystem,
		Ticker: ticker,
	}, nil
}


func (w *Watcher) Run(ctx context.Context) {
	defer w.Ticker.Stop()
	for {
		select{
		case <-ctx.Done():
			return
		case <-w.Ticker.C():
			currentSnapshot, err := w.FileSystem.Scan(ctx, w.Path)
			if err != nil {
				if errors.Is(err, os.ErrNotExist) && w.FolderState != FolderStateUnavailable{
					w.FolderState = FolderStateUnavailable
					w.Reporter.ReportFolderState(w.FolderState, w.Path)
				}
				continue
			}
			if w.FolderState == FolderStateUnavailable{
				w.FolderState = FolderStateAvailable
				w.Reporter.ReportFolderState(w.FolderState, w.Path)
			}
			changes := w.check(currentSnapshot)
			w.Reporter.Report(changes)
		}
	}
}

func (w *Watcher) check(currentSnapshot map[string]FileSnapshot) (changes Changes){
	changes.Added, changes.Removed, changes.Modified = compare(currentSnapshot, w.PreviousSnapshot)
	w.PreviousSnapshot = currentSnapshot
	return changes
}

func compare(currentSnapshot, previousSnapshot map[string]FileSnapshot) (added []string, removed []string, modified []string){
	for previousFile, previousFileContent := range previousSnapshot{
		currentFile, fileExist := currentSnapshot[previousFile]

		if !fileExist{
			removed = append(removed, previousFile)
		}else if currentFile.IsModified(previousFileContent){
			modified = append(modified, previousFile)
		}
	}
	for file := range currentSnapshot{
		_, fileExist := previousSnapshot[file]
		if !fileExist{
			added = append(added, file)
		}
	}
	return added, removed, modified
}