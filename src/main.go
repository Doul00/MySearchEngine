package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// Package global variable to get the file where the index is saved
var pathToIndexSave string

func main() {
	args := os.Args[1:]
	var pathToFile string

	if len(args) < 1 {
		fmt.Println("USAGE: ./main path_to_your_folder")
		fmt.Println("or ./main path_to_your_folder index_save_filename")
		panic("Please follow the usage")
	}

	if len(args) == 1 {
		pathToFile = args[0]
		pathToIndexSave = "index_save"
	} else if len(args) == 2 {
		pathToFile = args[0]
		pathToIndexSave = args[1]
	}

	dir, _ := filepath.Abs(pathToFile)
	documentsList := fetch(dir, true)
	processedDocs := processDocuments(documentsList[:100])
	postings := index(processedDocs)
	index := build(postings, processedDocs)
	save(index, pathToIndexSave)
	index = load(pathToIndexSave)
	pages := search("action and christ")
	fmt.Println(pages)

}
