package reporter

import (
	"fmt"

	"github.com/martin0b101/watcher/internal/model"
)

type ConsoleReporter struct {
}

// ReportFolderState implements Reporter.
func (ConsoleReporter) ReportMessage(msg string) {
	fmt.Println(msg)
}

func (cr ConsoleReporter) Report(changes model.Changes) {
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
