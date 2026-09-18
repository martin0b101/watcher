package watcher

import (
	"fmt"
	"time"
)


type WatcherTicker struct {
	ticker *time.Ticker
}

func NewWatcherTicker(interval time.Duration) (*WatcherTicker, error) {
	if interval <= 0 {
		return nil, fmt.Errorf("error invalid interval")
	}
	ticker := time.NewTicker(interval)
	return &WatcherTicker{
		ticker: ticker,
	}, nil
}

// C implements Ticker.
func (wt WatcherTicker) C() <-chan time.Time {
	return wt.ticker.C
}

// Stop implements Ticker.
func (wt WatcherTicker) Stop() {
	wt.ticker.Stop()
}
