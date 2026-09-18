package processor

import (
	"context"
	"fmt"
	"testing"

	"github.com/martin0b101/watcher/internal/model"
	"github.com/stretchr/testify/assert"
)

type ScenatioTest struct {
	WantError bool
	Err error
}

type MoveFiles struct {
	Source string
	Destination string
}

type MockFileSystem struct {
	Scenarios []ScenatioTest
	MovedFiles []MoveFiles
	moveCounter int
}

func (mfs *MockFileSystem) Scan(ctx context.Context, path string) (map[string]model.FileSnapshot, error) {
	return map[string]model.FileSnapshot{}, nil
}

func (mfs *MockFileSystem) Move(ctx context.Context, source string, dest string) error{
	scenarioResult := mfs.Scenarios[mfs.moveCounter]
	mfs.moveCounter += 1
	if scenarioResult.WantError{
		return scenarioResult.Err
	}
	mfs.MovedFiles = append(mfs.MovedFiles, MoveFiles{
		Source: source,
		Destination: dest,
	})
	return nil 
}

func TestProcess(t *testing.T) {
	ctx := context.Background()
	mockFileSystem := MockFileSystem{
		Scenarios: []ScenatioTest{
			{
				WantError: false,
				Err: nil,
			},
			{
				WantError: true,
				Err: fmt.Errorf("failed to move the file"),
			},
			{
				WantError: false,
				Err: nil,
			},
		},
	}

	testFiles := []model.FileSnapshot{
		{
			Name: "a.txt",
			FullPath: "test-folder/a.txt",
			Format: ".txt",
		},
{
			Name: "b.pdf",
			FullPath: "test-folder/a.pdf",
			Format: ".pdf",
		},
		{
			Name: "c.csv",
			FullPath: "test-folder/c.csv",
			Format: ".csv",
		},
	}


	processor, err := NewFileProcessor(&mockFileSystem)
	assert.NoError(t, err)
	processor.Process(ctx, testFiles)

	assert.Equal(t, "test-folder/documents/a.txt", mockFileSystem.MovedFiles[0].Destination)
	assert.Equal(t, "test-folder/data/c.csv", mockFileSystem.MovedFiles[1].Destination)
	assert.Len(t, mockFileSystem.MovedFiles, 2)
}