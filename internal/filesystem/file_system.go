package filesystem

import (
	"context"
	"os"
	"path/filepath"
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

type FileSystem interface {
	Scan(ctx context.Context, path string) (map[string]FileSnapshot, error)
}

type OSFileSystem struct {
	Recursive bool
}

var _ FileSystem = OSFileSystem{}

func (osfs OSFileSystem) Scan(ctx context.Context, path string) (map[string]FileSnapshot, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	return scanFolder(ctx, files, path, path, osfs.Recursive)
}

func scanFolder(ctx context.Context, files []os.DirEntry, basePath string, currentPath string, recursive bool) (map[string]FileSnapshot, error) {
	currentSnapshot := make(map[string]FileSnapshot)
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
		currentSnapshot[relPath] = FileSnapshot{
			Name:       relPath,
			ModifiedAt: info.ModTime(),
			Size:       info.Size(),
		}
	}
	return currentSnapshot, nil
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
