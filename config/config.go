package config

import (
	"errors"
	"flag"
	"fmt"
	"time"
)

type OutputFormat string

const (
	Console OutputFormat = "console"
	JSON    OutputFormat = "json"
)

type OutputFormatFlag struct {
	Value OutputFormat
}

func (off *OutputFormatFlag) String() string {
	return string(off.Value)
}

func (off *OutputFormatFlag) Set(val string) error {
	switch OutputFormat(val) {
	case Console, JSON:
		off.Value = OutputFormat(val)
		return nil
	default:
		return errors.New("must be one of: console, json")
	}
}

type Config struct {
	Path         string
	Interval     time.Duration
	Recursive    bool
	OutputFormat OutputFormat
}

func Parse(args []string) (*Config, error) {
	fs := flag.NewFlagSet("watcher", flag.ContinueOnError)
	outputFlag := OutputFormatFlag{
		Value: Console,
	}
	interval := fs.Duration("interval", 3 * time.Second, "scaning interval in seconds")
	recursive := fs.Bool("recursive", false, "scaning folders recursively")
	fs.Var(&outputFlag, "output", "output format (e.g. console, json)")

	if len(args) > 2{
		err := fs.Parse(args[1:])
	
		if err != nil{
			return nil, err
		}
	}
	if fs.NArg() <= 0 {
        return nil, fmt.Errorf("required arguments missing")
    }

	return &Config{
		Path:         fs.Arg(0),
		Interval:     *interval,
		Recursive:    *recursive,
		OutputFormat: outputFlag.Value,
	}, nil
}
