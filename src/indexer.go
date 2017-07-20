package main

import (
	"encoding/gob"
	"os"
)

/*
* Structs
 */

// Generation represents a version of the index
type Generation struct {
	WordsToDid map[string][]int
}

/* Index is the index of the search engine
* urlToDid is a mapping of each url to its did
* didToUrl is the inverse map of urlToDid, given a did we have the url
* generations is the list of all the index's generations
 */
type Index struct {
	TokDocs []TokenizedDocument
	Posting map[string][]string
}

/*
* Functions
 */
func buildIndex(docs []Document, index string) {
	tokDocs := processDocuments(&docs)
	posting := buildPosting(&tokDocs) // Map of words to list of urls

	save_index(&Index{tokDocs, posting}, index)
}

/*
* @documents the list of documents and the words they contains
* Returns a list of postings
 */
func buildPosting(documents *[]TokenizedDocument) map[string][]string {
	posting := make(map[string][]string)

	for _, doc := range *documents {
		// Prevents adding the same document several times.
		var flagInsertion = make(map[string]bool)

		for word := range doc.Title {
			_, isInserted := flagInsertion[word]

			if len(posting[word]) == 0 && !isInserted {
				posting[word] = []string{doc.Url}
			} else if !isInserted {
				posting[word] = append(posting[word], doc.Url)
			}
		}

		for word := range doc.Body {
			_, isInserted := flagInsertion[word]

			if len(posting[word]) == 0 && !isInserted {
				posting[word] = []string{doc.Url}
			} else if !isInserted {
				posting[word] = append(posting[word], doc.Url)
			}
		}
	}

	return posting
}

/*
* @index the index to save
* @path the file in which the index will be saved
* Saves the index into a file using the goland 'gob' serializer
 */
func save_index(index *Index, path string) {
	f, err := os.Create(path)
	defer f.Close()

	handleError(err, "Error saving the index at path:"+path)
	err = gob.NewEncoder(f).Encode(*index)
	handleError(err, "Error encoding the index:"+path)
}
