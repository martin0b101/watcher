package main

import (
	"fmt"
	"io"
	"os"
)

type Changes struct {
	Added []string
	Removed []string
	Modified []string
}

type Reporter interface {
    Report(changes Changes)
}

type ConsoleReporter struct {
	out io.Writer
}

var _ Reporter = ConsoleReporter{}

func NewConsoleReporter(out io.Writer) ConsoleReporter{
	if out == nil {
		out = os.Stdout
	}
	return ConsoleReporter{
		out: out,
	}
}

func (cr ConsoleReporter) Report(changes Changes){
	for _, file := range changes.Added{
		fmt.Fprintf(cr.out, "New file: %s\n", file)
	}
	for _, file := range changes.Removed{
		fmt.Fprintf(cr.out, "Removed: %s\n", file)
	}
	for _, file := range changes.Modified{
		fmt.Fprintf(cr.out, "Modified: %s\n", file)
	}
}
