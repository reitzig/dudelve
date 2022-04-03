package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit
		}

	case ScanProgress:
		if msg.scannedOk {
			m.scannedFiles += uint64(msg.scannedCount)
			return m, func() tea.Msg { return msg.awaitNext() }
		} else {
			// Scanning finished!
			err, ok := <-msg.err
			if ok {
				// TODO: display error nicely?
				log.Fatal(err)
				return m, tea.Quit
			}

			m.currentDir, _ = <-msg.result
		}
	}

	return m, nil
}
