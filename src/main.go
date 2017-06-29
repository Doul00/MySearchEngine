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
		fmt.Println("USAGE: To generate an index from your documents: \n ./main --create-index path_to_your_folder index_filename")
		fmt.Println("To launch the search engine: \n ./main --search index_filename \n Then type your request in the format described in the README")
		fmt.Println("For more information please consult the README")
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
	pages := search("christians")
	fmt.Println(pages)

}
