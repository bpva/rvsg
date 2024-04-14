package main

import (
	"fmt"
	"io"
	"os"
	"sort"
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
	return dirTreeRecursive(out, path, printFiles, 0, "")
}

func dirTreeRecursive(out io.Writer, path string, printFiles bool, level int, indentLine string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	fs, err := file.Stat()
	if err != nil {
		return err
	}
	if !fs.IsDir() {
		printName(out, fs, true)
		return nil
	}

	files, err := file.Readdir(0)
	if err != nil {
		return err
	}

	if !printFiles {
		files = filterDirectories(files)
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	for i, f := range files {
		isLast := i == len(files)-1
		if f.IsDir() {
			fmt.Fprint(out, indentLine)
			printName(out, f, isLast)
			newIndent := indentLine
			if isLast {
				newIndent += "\t"
			} else {
				newIndent += "│\t"
			}
			dirTreeRecursive(out, path+"/"+f.Name(), printFiles, level+1, newIndent)
		} else {
			fmt.Fprint(out, indentLine)
			printName(out, f, isLast)
		}
	}

	return nil
}

func printName(out io.Writer, file os.FileInfo, isLast bool) {
	if isLast {
		fmt.Fprint(out, lastBranch)
	} else {
		fmt.Fprint(out, branchChar)
	}

	if file.IsDir() {
		fmt.Fprintf(out, "───%s\n", file.Name())
		return
	} else {
		fileSize := file.Size()
		size := ""
		if fileSize == 0 {
			size = "empty"
		} else {
			size = fmt.Sprintf("%db", fileSize)
		}
		fmt.Fprintf(out, "───%s (%s)\n", file.Name(), size)
	}
}

func filterDirectories(files []os.FileInfo) []os.FileInfo {
	dirs := make([]os.FileInfo, 0)
	for _, f := range files {
		if f.IsDir() {
			dirs = append(dirs, f)
		}
	}
	return dirs
}
