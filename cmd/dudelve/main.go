package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// TODO: CLI interface!
var rootDir string
var showHidden bool

//var backupExclusionTags []string

func main() {
	rootDir = "/home/raphael/Download"
	showHidden = false

	// TODO: set up logger

	p := tea.NewProgram(initialModel(rootDir), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func (m model) Init() tea.Cmd {
	progressOut := make(chan uint64, 100)
	errorOut := make(chan error, 1)
	resultOut := make(chan *node, 1)

	go func() {
		defer close(progressOut)
		defer close(errorOut)
		defer close(resultOut)

		rootNode, err := scanTree(m.currentPath, nil, progressOut)
		if err != nil {
			errorOut <- err
		} else {
			resultOut <- rootNode
		}
	}()

	return func() tea.Msg {
		return ScanProgress{moreCounts: progressOut, err: errorOut, result: resultOut}.awaitNext()
	}
}

// NB: On local filesystems at least, this is plenty fast. Caching or parent pointers can come later.
func scanTree(filePath string, parent *node, progress chan uint64) (*node, error) {
	// TODO: capture permission errors; but what to do? Warn and ignore?
	fileInfo, err := os.Lstat(filePath)
	if err != nil {
		return nil, err
	}

	currentNode := node{name: fileInfo.Name(), parent: parent, count: 1}

	if fileInfo.IsDir() {
		accumulatedSize := uint64(0)
		count := uint64(0)
		children := []*node{}

		files, err := ioutil.ReadDir(filePath)
		if err != nil {
			return nil, err
		}

		for _, file := range files {
			if !showHidden && strings.HasPrefix(file.Name(), ".") {
				continue
			}

			child, err := scanTree(filepath.Join(filePath, file.Name()), &currentNode, progress)
			if err != nil {
				return nil, err
			}
			accumulatedSize += child.size
			count += child.count
			children = append(children, child)
		}

		currentNode.size = accumulatedSize
		currentNode.count += count
		currentNode.children = children
		progress <- uint64(len(children))
		//time.Sleep(50 * time.Millisecond)
	} else if fileInfo.Mode()&os.ModeSymlink != 0 { // symlink
		return &node{name: fileInfo.Name()}, nil
	} else { // "normal" file
		currentNode.size = uint64(fileInfo.Size())
	}

	return &currentNode, nil
}
