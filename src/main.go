package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Package global variable to get the file where the index is saved
var pathToIndexSave *string

func main() {
	var idx Index
	var postings []Posting

	mode := flag.String("mode", "create", "the program mode (Required)")
	pathToFolder := flag.String("path", "", "the path to your folder containing the files (Required)")
	pathToIndexSave = flag.String("index", "", "the name of the index file (Required)")

	flag.Parse()

	// Launch the indexing
	if strings.Compare(*mode, "create") == 0 {

		if len(*pathToFolder) == 0 || len(*pathToIndexSave) == 0 {
			usage("Please follow the usage")
		}

		// Checks if the index already exists
		f, err := os.Open(*pathToIndexSave)
		defer f.Close()

		dir, _ := filepath.Abs(*pathToFolder)

		fmt.Println("Reading the documents...")
		documentsList := fetch(dir, true)

		fmt.Println("Processing the documents...")
		processedDocs := processDocuments(documentsList[:100])

		postings = index(processedDocs)

		if err == nil {
			fmt.Println("Updating the index...")
			idx = load(*pathToIndexSave)
			updateGeneration(postings, &idx)
		} else {
			fmt.Println("Creating the index...")
			idx = build(postings, processedDocs)
		}

		fmt.Println("Saving the index in file: " + *pathToIndexSave)
		save(idx, *pathToIndexSave)
	}

	// Launch the search prompt
	if strings.Compare(*mode, "search") == 0 {
		if len(*pathToIndexSave) == 0 {
			usage("Please follow the usage.")
		}

		index := load(*pathToIndexSave)

		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Type exit to quit the shell")
		fmt.Println("---------------------")

		for {
			fmt.Print("> ")
			query, _ := reader.ReadString('\n')
			query = strings.Replace(query, "\n", "", -1)
			if strings.Compare(query, "exit") == 0 {
				os.Exit(0)
			}
			pages := search(query, index)
			formatAnswers(pages)
		}

	}

}

func usage(message string) {
	fmt.Println("\n\nUSAGE: To generate an index from your documents: \n./main -mode=create -path=path_to_your_folder -index=index_filename")
	fmt.Println("To launch the search engine: \n./main -mode=search -index=index_filename \nThen type your request in the format described in the README")
	fmt.Println("For more information please consult the README")
	panic(message)
}

func formatAnswers(results []string) {
	fmt.Println("\n" + strconv.Itoa(len(results)) + " result(s) found")
	for i := range results {
		fmt.Println(results[i])
	}
}
