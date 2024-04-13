package main

import (
	"io"
	"os"
)

const (
	branchChar = "├"
	indentChar = "│"
	lineChar   = "─"
	lastBranch = "└"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	return dirTreeRecursive(out, path, printFiles, 0)
}
