package main

import (
	"fmt"
	"time"
)

type Watcher struct {
	Path string
	PreviousSnapshot map[string]FileSnapshot
	Interval time.Duration
	Reporter Reporter
	FileSystem FileSystem
}

func NewWatcher(path string, interval time.Duration, reporter Reporter, fileSystem FileSystem) (*Watcher, error) {
	if path == ""{
		return nil, fmt.Errorf("path is empty")
	}

	if reporter == nil{
		return nil, fmt.Errorf("reporter is not init")
	}

	if fileSystem == nil {
		return nil, fmt.Errorf("fileSystem is not init")
	}

	snapshot, err := fileSystem.Scan(path)
	if err != nil{
		return nil, fmt.Errorf("failed to initialize snapshot: %w", err)
	}

	return &Watcher{
		Path: path,
		PreviousSnapshot: snapshot,
		Interval: interval,
		Reporter: reporter,
		FileSystem: fileSystem,
	}, nil
}


func (w *Watcher) Run() {
	for {
		time.Sleep(w.Interval)
		currentSnapshot, err := w.FileSystem.Scan(w.Path)
		if err != nil{
			fmt.Printf("error scaning: %v\n", err)
			continue
		}
		changes := w.check(currentSnapshot)
		w.Reporter.Report(changes)
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