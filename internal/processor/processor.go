package processor

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/martin0b101/watcher/internal/model"
	"github.com/martin0b101/watcher/internal/watcher"
)


type destination string

const (
	destinationDocument = "documents"
	destinaitionImage = "images"
	destinationData = "data"
)

func (d destination) String() string{
	return string(d)
}

var exstantionMatcher = map[string]destination{
	".pdf": destinationDocument,
	".txt": destinationDocument,
	".jpg": destinaitionImage,
	".csv": destinationData,
}

type FileProcessor struct {
	fileSystem watcher.FileSystem
}

func NewFileProcessor(system watcher.FileSystem) (*FileProcessor, error){
	if system == nil{
		return nil, fmt.Errorf("can not init file system")
	}
	return &FileProcessor{
		fileSystem: system,
	}, nil
}

func (fp FileProcessor) Process(ctx context.Context, files []model.FileSnapshot){
	for _, file := range files{
		if ctx.Err() !=nil {
			return
		}
		destinationFolder, ok := exstantionMatcher[file.Format]
		if !ok{
			continue
		}
		dest := destinationFolder.String()
		folder := filepath.Dir(file.FullPath)
		destPath := filepath.Join(folder, dest, file.Name)
		err := fp.fileSystem.Move(ctx, file.FullPath, destPath)
		if err != nil{
			fmt.Printf("Failed to process %s: %v\n", file.Name, err)
			continue
		}
		fmt.Printf("File processed: %s\n", file.Name)
	}
}