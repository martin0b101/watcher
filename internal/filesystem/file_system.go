package filesystem

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/martin0b101/watcher/internal/model"
)

var ErrFolderUnavailable = errors.New("folder unavailable")

type OSFileSystem struct {
	Recursive bool
}

func (osfs OSFileSystem) Scan(ctx context.Context, path string) (map[string]model.FileSnapshot, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	files, err := os.ReadDir(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist){
			return nil, ErrFolderUnavailable
		}
		return nil, err
	}
	return scanFolder(ctx, files, path, path, osfs.Recursive)
}

func scanFolder(ctx context.Context, files []os.DirEntry, basePath string, currentPath string, recursive bool) (map[string]model.FileSnapshot, error) {
	currentSnapshot := make(map[string]model.FileSnapshot)
	for _, file := range files {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		fullPath := filepath.Join(currentPath, file.Name())

		if file.IsDir() {
			if !recursive {
				continue
			}
			subFiles, err := os.ReadDir(fullPath)
			if err != nil {
				if errors.Is(err, os.ErrNotExist){
					return nil, ErrFolderUnavailable
				}
				return nil, err
			}
			newCurrentSnapshot, err := scanFolder(ctx, subFiles, basePath, fullPath, true)
			if err != nil {
				return nil, err
			}
			currentSnapshot = MergeMaps(currentSnapshot, newCurrentSnapshot)
			continue
		}
		info, err := file.Info()
		if err != nil {
			return nil, err
		}

		relPath, err := filepath.Rel(basePath, fullPath)
		if err != nil {
			return nil, err
		}
		currentSnapshot[relPath] = model.FileSnapshot{
			Name:       relPath,
			ModifiedAt: info.ModTime(),
			Size:       info.Size(),
			Format:     filepath.Ext(info.Name()),
			FullPath: fullPath,
		}
	}
	return currentSnapshot, nil
}


func (osfs OSFileSystem) Move(ctx context.Context,sourcePath string, destinationPath string) error {
	if ctx.Err() != nil{
		return ctx.Err()
	}
	
	destDir := filepath.Dir(destinationPath)
	
	err := os.MkdirAll(destDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}
	
	err = os.Rename(sourcePath, destinationPath)
	if err != nil {
		return err 
	}

	return nil
}

func MergeMaps[K comparable, V any](map1, map2 map[K]V) map[K]V {
	merged := make(map[K]V, len(map1)+len(map2))

	for k, v := range map1 {
		merged[k] = v
	}
	for k, v := range map2 {
		merged[k] = v
	}

	return merged
}
