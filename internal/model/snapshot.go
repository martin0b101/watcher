package model

import "time"

type FileSnapshot struct {
	Name       string
	ModifiedAt time.Time
	Size       int64
}

// because the method doesn't modify f, so there's no reason to use a pointer receiver.
func (f FileSnapshot) IsModified(other FileSnapshot) bool {
	return f.ModifiedAt.After(other.ModifiedAt) || f.Size != other.Size
}