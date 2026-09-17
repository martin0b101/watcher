package main

import (
	"fmt"
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
}

var _ Reporter = ConsoleReporter{}


func (cr ConsoleReporter) Report(changes Changes){
	for _, file := range changes.Added{
		fmt.Printf("New file: %s\n", file)
	}
	for _, file := range changes.Removed{
		fmt.Printf("Removed: %s\n", file)
	}
	for _, file := range changes.Modified{
		fmt.Printf("Modified: %s\n", file)
	}
}
