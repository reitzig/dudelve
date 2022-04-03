package main

import (
	"fmt"
	"math/big"

	"github.com/dustin/go-humanize"
)

func (m model) View() string {
	view := fmt.Sprintf("Path: %s\nScanned: %d", m.currentPath, m.scannedFiles)
	if m.currentDir != nil {
		view += fmt.Sprintf("\n%s: %d files, %s",
			m.currentDir.name,
			m.currentDir.count,
			humanize.BigIBytes(big.NewInt(int64(m.currentDir.size))))
	}

	return view
}
