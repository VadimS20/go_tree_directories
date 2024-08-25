package main

import (
	"fmt"
	"io/fs"
	"os"
)

func remove(slice []fs.FileInfo, s int) []fs.FileInfo {
	return append(slice[:s], slice[s+1:]...)
}

func helper(out *os.File, path string, flag bool, tempPath string) error {
	dir, _ := os.Open(path)
	defer dir.Close()

	files, _ := dir.Readdir(-1)

	if !flag {
		for i := 0; i < len(files); i++ {
			if !files[i].IsDir() {
				files = remove(files, i)
			}
		}
	}

	for i, file := range files {
		if file.IsDir() {
			tabs := tempPath
			if i == len(files)-1 {
				tabs = tabs + "\t"
				answer := tempPath + "└───" + file.Name()
				fmt.Fprintln(out, answer)
			} else {
				tabs = tabs + "│	"
				answer := tempPath + "├───" + file.Name()
				fmt.Fprintln(out, answer)
			}
			newPath := path + "\\" + file.Name()
			helper(out, newPath, flag, tabs)
		} else {
			if flag {
				if i == len(files)-1 {
					answer := tempPath + "└───" + file.Name()
					fmt.Fprintln(out, answer)
				} else {
					answer := tempPath + "├───" + file.Name()
					fmt.Fprintln(out, answer)
				}
			}
		}
	}
	return nil
}

func dirTree(out *os.File, path string, flag bool) error {
	return helper(out, path, flag, "")
}

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
