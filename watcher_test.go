package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCompateEmpty(t *testing.T) {
	previous := map[string]FileSnapshot{
		"a.txt": {},
		"b.txt": {},
	}

	current := map[string]FileSnapshot{
		"a.txt": {},
		"b.txt": {},
	}

	added, removed, modified := compare(current, previous)

	assert.Empty(t, added)
	assert.Empty(t, removed)
	assert.Empty(t, modified)
}

func TestCompareAddedFile(t *testing.T) {
	previous := map[string]FileSnapshot{
		"a.txt": {},
	}

	current := map[string]FileSnapshot{
		"a.txt": {},
		"b.txt": {},
	}

	added, removed, modified := compare(current, previous)

	assert.Len(t, added, 1)
	assert.Contains(t, added, "b.txt")
	assert.Empty(t, removed)
	assert.Empty(t, modified)
}

func TestCompareRemovedFile(t *testing.T) {
	previous := map[string]FileSnapshot{
		"a.txt": {},
		"b.txt": {},
	}

	current := map[string]FileSnapshot{
		"a.txt": {},
	}

	added, removed, modified := compare(current, previous)

	assert.Len(t, removed, 1)
	assert.Contains(t, removed, "b.txt")
	assert.Empty(t, added)
	assert.Empty(t, modified)
}

func TestCompareModifiedFile_ModifiedAt(t *testing.T) {
	previous := map[string]FileSnapshot{
		"a.txt": {},
		"b.txt": {
			ModifiedAt: time.Date(2000, 1, 1, 11, 45, 0, 0, time.UTC),
		},
	}

	current := map[string]FileSnapshot{
		"a.txt": {},
		"b.txt": {
			ModifiedAt: time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC),
		},
	}

	added, removed, modified := compare(current, previous)

	assert.Len(t, modified, 1)
	assert.Contains(t, modified, "b.txt")
	assert.Empty(t, removed)
	assert.Empty(t, added)
}

func TestCompareModifiedFile_Size(t *testing.T) {
	previous := map[string]FileSnapshot{
		"a.txt": {},
		"b.txt": {
			Size: 100,
		},
	}

	current := map[string]FileSnapshot{
		"a.txt": {},
		"b.txt": {
			Size: 200,
		},
	}

	added, removed, modified := compare(current, previous)

	assert.Len(t, modified, 1)
	assert.Contains(t, modified, "b.txt")
	assert.Empty(t, removed)
	assert.Empty(t, added)
}

type ScanResult struct {
	Snapshot map[string]FileSnapshot
	Err      error
}

type MockFileSystem struct {
	Snapshots []ScanResult
	scanCount int
}

func (mfs *MockFileSystem) Scan(path string) (map[string]FileSnapshot, error) {
	if mfs.scanCount >= len(mfs.Snapshots) {
		return nil, fmt.Errorf("no more snapshots available")
	}

	snapshotToReturn := mfs.Snapshots[mfs.scanCount]

	mfs.scanCount += 1
	if snapshotToReturn.Snapshot != nil {
		return snapshotToReturn.Snapshot, nil
	}
	return nil, snapshotToReturn.Err
}

type MockTicker struct {
	channel chan time.Time
}

func (mt MockTicker) C() <-chan time.Time {
	return mt.channel
}

// Stop implements Ticker.
func (mt MockTicker) Stop() {
}

type MockReport struct {
	Changes Changes
}

func (mr *MockReport) Report(changes Changes) {
	mr.Changes = changes
}

func (mr *MockReport) ReportFolderState(FolderState FolderState, path string) {

}

func TestChange(t *testing.T) {
	testReporter := ConsoleReporter{}
	fakeTicker, err := NewWatcherTicker(time.Minute)
	assert.NoError(t, err)
	mockFileSystem := MockFileSystem{
		Snapshots: []ScanResult{
			{
				Snapshot: map[string]FileSnapshot{

					"a.txt": {
						Name: "a.txt",
					},
					"c.txt": {
						Name: "c.txt",
					},
					"d.txt": {
						Name: "d.txt",
						Size: 100,
					},
				},
			},
		},
	}

	currentSnapshot := map[string]FileSnapshot{
		"a.txt": {
			Name: "a.txt",
		},
		"b.txt": {
			Name: "b.txt",
		},
		"d.txt": {
			Name: "d.txt",
			Size: 200,
		},
	}

	testWatcher, err := NewWatcher("test-path", testReporter, &mockFileSystem, fakeTicker)
	assert.NoError(t, err)

	changes := testWatcher.check(currentSnapshot)

	assert.Len(t, changes.Added, 1)
	assert.Len(t, changes.Removed, 1)
	assert.Len(t, changes.Modified, 1)
	assert.Contains(t, changes.Added, "b.txt")
	assert.Contains(t, changes.Removed, "c.txt")
	assert.Contains(t, changes.Modified, "d.txt")

	assert.Equal(t, currentSnapshot, testWatcher.PreviousSnapshot)
}

func TestWatchContrutor(t *testing.T) {
	testReporter := ConsoleReporter{}
	fakeTicker, err := NewWatcherTicker(time.Minute)
	assert.NoError(t, err)
	mockFileSystem := MockFileSystem{
		Snapshots: []ScanResult{
			{
				Snapshot: map[string]FileSnapshot{
					"a.txt": {
						Name: "a.txt",
					},
					"c.txt": {
						Name: "c.txt",
					},
					"d.txt": {
						Name: "d.txt",
						Size: 100,
					},
				},
			},
		},
	}

	expectedInitalState := map[string]FileSnapshot{
		"a.txt": {
			Name: "a.txt",
		},
		"c.txt": {
			Name: "c.txt",
		},
		"d.txt": {
			Name: "d.txt",
			Size: 100,
		},
	}

	testWatcher, err := NewWatcher("test-path", testReporter, &mockFileSystem, fakeTicker)
	assert.NoError(t, err)
	assert.Equal(t, testWatcher.PreviousSnapshot, expectedInitalState)
}

func TestWatchContrutor_ScanError(t *testing.T) {
	testReporter := ConsoleReporter{}
	fakeTicker, err := NewWatcherTicker(time.Minute)
	assert.NoError(t, err)
	mockFileSystem := MockFileSystem{
		Snapshots: []ScanResult{
			{
				Snapshot: map[string]FileSnapshot{
					"a.txt": {
						Name: "a.txt",
					},
					"c.txt": {
						Name: "c.txt",
					},
					"d.txt": {
						Name: "d.txt",
						Size: 100,
					},
				},
			},
		},
	}

	testWatcher, err := NewWatcher("test-path", testReporter, &mockFileSystem, fakeTicker)
	assert.Error(t, err)
	assert.Nil(t, testWatcher)
}

func TestWatcherRun(t *testing.T) {
	mockTicker := MockTicker{
		channel: make(chan time.Time),
	}
	testReporter := MockReport{}
	mockFileSystem := MockFileSystem{
		Snapshots: []ScanResult{
			{
				Snapshot: map[string]FileSnapshot{
					"a.txt": {
						Name: "a.txt",
					},
					"c.txt": {
						Name: "c.txt",
					},
					"d.txt": {
						Name: "d.txt",
						Size: 100,
					},
				},
			}, {
				Snapshot: map[string]FileSnapshot{
					"a.txt": {
						Name: "a.txt",
					},
					"b.txt": {
						Name: "b.txt",
					},
					"d.txt": {
						Name: "d.txt",
						Size: 200,
					},
				},
			},
		},
	}
	testWatcher, err := NewWatcher("test-path", &testReporter, &mockFileSystem, mockTicker)
	assert.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan struct{})
	go func() {
		testWatcher.Run(ctx)
		done <- struct{}{}
	}()

	mockTicker.channel <- time.Now()

	wantedSnapshot := map[string]FileSnapshot{
		"a.txt": {
			Name: "a.txt",
		},
		"b.txt": {
			Name: "b.txt",
		},
		"d.txt": {
			Name: "d.txt",
			Size: 200,
		},
	}
	cancel()
	<-done

	assert.Equal(t, testWatcher.PreviousSnapshot, wantedSnapshot)
	assert.Len(t, testReporter.Changes.Added, 1)
	assert.Len(t, testReporter.Changes.Removed, 1)
	assert.Len(t, testReporter.Changes.Modified, 1)
}

func TestWatcherScanErrorRecovery(t *testing.T) {
	mockTicker := MockTicker{
		channel: make(chan time.Time),
	}
	testReporter := MockReport{}
	mockFileSystem := MockFileSystem{
		Snapshots: []ScanResult{
			{
				Snapshot: map[string]FileSnapshot{
					"a.txt": {
						Name: "a.txt",
					},
					"c.txt": {
						Name: "c.txt",
					},
					"d.txt": {
						Name: "d.txt",
						Size: 100,
					},
				},
			},
			{
				Err: fmt.Errorf("error while scaining"),
			},
			{
				Snapshot: map[string]FileSnapshot{
					"a.txt": {
						Name: "a.txt",
					},
					"b.txt": {
						Name: "b.txt",
					},
					"d.txt": {
						Name: "d.txt",
						Size: 200,
					},
				},
			},
		},
	}

	testWatcher, err := NewWatcher("test-path", &testReporter, &mockFileSystem, mockTicker)
	assert.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan struct{})
	go func() {
		testWatcher.Run(ctx)
		done <- struct{}{}
	}()

	mockTicker.channel <- time.Now()
	mockTicker.channel <- time.Now()

	wantedSnapshot := map[string]FileSnapshot{
		"a.txt": {
			Name: "a.txt",
		},
		"b.txt": {
			Name: "b.txt",
		},
		"d.txt": {
			Name: "d.txt",
			Size: 200,
		},
	}
	cancel()
	<-done

	assert.Equal(t, testWatcher.PreviousSnapshot, wantedSnapshot)
	assert.Len(t, testReporter.Changes.Added, 1)
	assert.Len(t, testReporter.Changes.Removed, 1)
	assert.Len(t, testReporter.Changes.Modified, 1)
}
