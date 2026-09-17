package main

import (
	"fmt"
	"time"
)

type Ticker interface {
	C() <-chan time.Time
	Stop()
}

var _ Ticker = WatcherTicker{}

type WatcherTicker struct {
	Ticker *time.Ticker
}

func NewWatcherTicker(interval time.Duration) (*WatcherTicker, error){
	if interval <= 0{
		return nil, fmt.Errorf("error invalid interval")
	}
	ticker := time.NewTicker(interval)
	return &WatcherTicker{
		Ticker: ticker,
	}, nil
}
// C implements Ticker.
func (wt WatcherTicker) C() <-chan time.Time {
	return wt.Ticker.C 
}

// Stop implements Ticker.
func (wt WatcherTicker) Stop() {
	wt.Ticker.Stop()
}
