package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCompateEmpty(t *testing.T){
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

func TestCompareRemovedFile(t *testing.T){
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

func TestCompareModifiedFile_ModifiedAt(t *testing.T){
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

func TestCompareModifiedFile_Size(t *testing.T){
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