package main

type model struct {
	currentPath  string
	scannedFiles uint64
	currentDir   *node
}

func initialModel(rootPath string) model {
	return model{currentPath: rootPath, scannedFiles: 0, currentDir: nil}
}

type node struct {
	name     string
	size     uint64
	children []*node
	count    uint64
	parent   *node
}
