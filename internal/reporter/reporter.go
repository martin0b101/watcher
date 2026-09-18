package reporter

import (
	"fmt"
)

type Changes struct {
	Added    []string
	Removed  []string
	Modified []string
}

type Reporter interface {
	Report(changes Changes)
	ReportMessage(msg string)
}

type ConsoleReporter struct {
}

// ReportFolderState implements Reporter.
func (ConsoleReporter) ReportMessage(msg string) {
	fmt.Println(msg)
}

var _ Reporter = ConsoleReporter{}

func (cr ConsoleReporter) Report(changes Changes) {
	for _, file := range changes.Added {
		fmt.Printf("New file: %s\n", file)
	}
	for _, file := range changes.Removed {
		fmt.Printf("Removed: %s\n", file)
	}
	for _, file := range changes.Modified {
		fmt.Printf("Modified: %s\n", file)
	}
}
