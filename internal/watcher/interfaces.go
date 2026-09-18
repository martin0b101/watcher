package watcher

import (
	"context"
	"time"

	"github.com/martin0b101/watcher/internal/model"
)

type FileSystem interface {
	Scan(ctx context.Context, path string) (map[string]model.FileSnapshot, error)
}

type Ticker interface {
	C() <-chan time.Time
	Stop()
}

type Reporter interface {
	Report(changes model.Changes)
	ReportMessage(msg string)
}