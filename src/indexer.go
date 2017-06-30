package main

import (
	"encoding/gob"
	"os"
)

/*
* Structs
 */

// Posting maps a word to the list of documents containing it
type Posting struct {
	word string
	url  []string
}

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
	UrlToDid    map[string]int
	DidToUrl    map[int]string
	Generations []Generation
}

/*
* Functions
 */

/*
* @documents the list of documents and the words they contains
* Returns a list of postings
 */
func index(documents []TokenizedDocument) []Posting {
	var postings []Posting
	var strMap = make(map[string][]string)

	for _, doc := range documents {
		// Prevents adding the same document twice
		var insertionMap = make(map[string]bool)

		for _, word := range doc.words {
			_, prs := strMap[word]
			_, isInserted := insertionMap[word]

			if !prs && !isInserted {
				newUrls := []string{doc.url}
				strMap[word] = newUrls
			} else if prs && !isInserted {
				strMap[word] = append(strMap[word], doc.url)
			}
			insertionMap[word] = true
		}
	}

	for k, v := range strMap {
		newPosting := Posting{word: k, url: v}
		postings = append(postings, newPosting)
	}

	return postings
}

/*
* @docs the list of TokenizedDocuments
* Returns a map containing, for each url, the corresponding did
 */
func createDid(docs []TokenizedDocument) (map[string]int, map[int]string) {
	m := make(map[string]int)
	inversedM := make(map[int]string)
	did := 0

	for _, doc := range docs {
		m[doc.url] = did
		inversedM[did] = doc.url
		did++
	}
	return m, inversedM
}

/*
* @didMap the link of every url to its did
* @urls the urls to convert to a list of the corresponding dids
* Returns the list of dids corresponding to the urls
 */
func urlsToDids(didMap map[string]int, urls []string) []int {
	var res = make([]int, len(urls))

	for i, str := range urls {
		res[i] = didMap[str]
	}
	return res
}

/*
* @didMap the link of every url to its did
* @postings the list of postings
* Returns a map containing a list of document dids for each word
 */
func createWordsToDid(didMap map[string]int, postings []Posting) map[string][]int {
	result := make(map[string][]int)

	for _, posting := range postings {
		result[posting.word] = urlsToDids(didMap, posting.url)
	}
	return result
}

func updateGeneration(postings []Posting, index *Index) {
	wordsToDid := createWordsToDid((*index).UrlToDid, postings)
	newGen := Generation{WordsToDid: wordsToDid}
	(*index).Generations = append((*index).Generations, newGen)
}

/*
* @postings the list of postings
* @docs the list of tokenized documents
* Creates an index containing the map of urls to dids and
* a map of words and the matching dids
 */
func build(postings []Posting, docs []TokenizedDocument) Index {
	urlToDID, didToURL := createDid(docs)
	wordsToDid := createWordsToDid(urlToDID, postings)
	newGen := Generation{WordsToDid: wordsToDid}
	return Index{UrlToDid: urlToDID, DidToUrl: didToURL, Generations: []Generation{newGen}}
}

/*
* @index the index to save
* @path the file in which the index will be saved
* Saves the index into a file using the goland 'gob' serializer
 */
func save(index Index, path string) {
	f, err := os.Create(path)

	handleError(err, "Error saving the index at path:"+path)
	enc := gob.NewEncoder(f)
	err = enc.Encode(index)
	handleError(err, "Error encoding the index:"+path)
	f.Close()
}
