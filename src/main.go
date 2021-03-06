package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

func main() {
	indexCommands := flag.NewFlagSet("index", flag.ExitOnError)
	searchCommands := flag.NewFlagSet("search", flag.ExitOnError)

	iData := indexCommands.String("path", "", "Path to the data to index.")
	iIndex := indexCommands.String("index", "index", "Prefix names for the indexes.")
	iNbGo := indexCommands.Int("go", 5, "Number of goroutines used to index the data.")
	iMaxDoc := indexCommands.Int("max", -1, "Maximum number of documents to index in the data.")

	sIndex := searchCommands.String("index", "index", "Path of index's folder.")
	sNbGo := searchCommands.Int("go", 5, "Number of goroutines used to search throught the indexes.")

	if len(os.Args) <= 1 {
		flag.PrintDefaults()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "index":
		indexCommands.Parse(os.Args[2:])
		indexCreation(*iData, *iIndex, *iNbGo, *iMaxDoc)

	case "search":
		searchCommands.Parse(os.Args[2:])
		searchEngine(*sIndex, *sNbGo)

	default:
		fmt.Println("Indexeur: 'index'")
		indexCommands.PrintDefaults()
		fmt.Println("Searcher: 'search'")
		searchCommands.PrintDefaults()
		os.Exit(2)
	}
}

func indexCreation(dataPath, index string, nbGo, maxDoc int) {
	data, err := os.Open(dataPath)
	if err != nil {
		panic(err)
	}
	defer data.Close()

	/* Creates Input channels for workers */
	docChans := make([]chan string, nbGo)
	for i := range docChans {
		docChans[i] = make(chan string)
	}

	/* Creates Stop channels for workers */
	stopChans := make([]chan bool, nbGo)
	for i := range stopChans {
		stopChans[i] = make(chan bool, 1)
	}

	/* Creates finish channels for workers */
	var wg sync.WaitGroup

	/* Creates all workers */
	for i := 0; i < nbGo; i++ {
		wg.Add(1)
		indexPath := index + "_" + strconv.Itoa(i) + ".idx"
		go workerIndex(docChans[i], stopChans[i], &wg, indexPath)
	}

	scanner := bufio.NewScanner(data)
	for i := 0; scanner.Scan(); i++ {
		if maxDoc != -1 && i >= maxDoc {
			break
		}
		docChans[i%nbGo] <- scanner.Text()
	}

	/* Stop workers' aggregation and starts indexation. */
	for i := range stopChans {
		stopChans[i] <- true
	}

	wg.Wait()
}

func workerIndex(docChan <-chan string, stopChan <-chan bool, wg *sync.WaitGroup, index string) {
	var docs []Document
	isDone := false

	for true {
		select {
		case <-stopChan:
			isDone = true
			break
		case rawDoc := <-docChan:
			docs = append(docs, createDocument(rawDoc))
		}

		if isDone {
			break
		}
	}

	buildIndex(docs, index)
	wg.Done()
}

func loadAllIndex(pattern string) []*Index {
	matches, err := filepath.Glob(pattern)
	if err != nil || len(matches) == 0 {
		panic(err)
	}

	var indexes []*Index
	for _, match := range matches {
		indexes = append(indexes, loadIndex(match))
	}

	return indexes
}

func searchEngine(indexPath string, nbGo int) {
	indexes := loadAllIndex(indexPath)

	scanner := bufio.NewScanner(os.Stdin)

	startDispatching(scanner, indexes, nbGo)
}

func usage(message string) {
	fmt.Println("\n\nUSAGE: To generate an index from your documents: \n./main -mode=create -path=path_to_your_folder -index=index_filename")
	fmt.Println("To launch the search engine: \n./main -mode=search -index=index_filename \nThen type your request in the format described in the README")
	fmt.Println("For more information please consult the README")
	panic(message)
}

func printResults(ans []TokenizedDocument) {
	for _, doc := range ans {
		fmt.Println(doc.Url)
	}
}
