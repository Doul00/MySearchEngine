package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"unicode"

	"io"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

/*
* Structs
 */

// Document contains the information regarding a document
type Document struct {
	text, url string
}

// TokenizedDocument contains the document and all its words
type TokenizedDocument struct {
	words []string
	url   string
}

// TextProcessor is a interface for all processors
type TextProcessor interface {
	process(str *string)
}

// DownCaseProcessor turns words into downcase
type DownCaseProcessor struct {
}

// AccentProcessor removes the accents from the words
type AccentProcessor struct {
}

/*
* Functions
 */

/*
* @error the error
* @msg a message describing the error
* Exits the program if e is not nil
 */
func handleError(e error, msg string) {
	if e != nil {
		fmt.Println(msg)
		panic(e)
	}
}

/*
* @input the bytes array to convert to UTF8
* Returns a string converted to UTF8
 */
func toUtf8(input []byte) string {
	buf := make([]rune, len(input))
	for i, b := range input {
		buf[i] = rune(b)
	}
	return string(buf)
}

/*
* @path is the path to the file to read
* Returns a document contenaining the text read
 */
func readFile(path string) Document {
	text, err := ioutil.ReadFile(path)
	handleError(err, "Cannot read file at path: "+path)
	doc := Document{url: path, text: toUtf8(text)}
	return doc
}

/*
* @path is the path to the file or to the folder
* Returns an array of pointers to the documents created
 */
func getDirDocuments(path string) []*Document {

	var result []*Document

	files, err := ioutil.ReadDir(path)
	handleError(err, "Cannot read directory at path: "+path)

	for _, f := range files {
		var tmpList []*Document

		if f.IsDir() {
			tmpList = getDirDocuments(path + "/" + f.Name())
		} else {
			doc := readFile(path + "/" + f.Name())
			tmpList = append(tmpList, &doc)
		}

		result = append(result, tmpList...)
	}
	return result
}

/*
* @recursive if true, fetch goes into directories
* @path is the path to the folder
* Returns a slice containing pointers to documents
 */
func fetch(path string, recursive bool) []*Document {

	var result []*Document

	// Open the directory and iterates through the files
	files, err := ioutil.ReadDir(path)
	handleError(err, "Cannot open directory at: "+path)

	for _, f := range files {
		var tmpList []*Document
		var doc Document

		if f.IsDir() && recursive {
			tmpList = getDirDocuments(path + "/" + f.Name())
		} else if !f.IsDir() {
			// Skips OSX hidden files (.DS_Store)
			if f.Name()[0] != '.' {
				doc = readFile(path + "/" + f.Name())
			} else {
				continue
			}
		}

		if len(tmpList) > 0 {
			result = append(result, tmpList...)
		} else if &doc != nil {
			result = append(result, &doc)
		}
	}
	return result
}

/*
* @str the string to process
* Removes the accents from the string
 */
func (p AccentProcessor) process(str *string) {
	b := make([]byte, len(*str))
	var r io.Reader = strings.NewReader(*str)

	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	r = transform.NewReader(r, t)
	io.ReadFull(r, b)
	*str = string(b)
}

/*
* @str the string to process
* Puts the string in downcase
 */
func (p DownCaseProcessor) process(str *string) {
	*str = strings.ToLower(*str)
}

/*
* @document the document to process
* @processors the processors transforming the data
* Returns a document which text has been downcased and cleaned (no accents and other symbols)
 */
func analyse(document Document, processors []TextProcessor) TokenizedDocument {

	processedText := document.text
	for _, p := range processors {
		p.process(&processedText)
	}

	re := regexp.MustCompile("[[:^word:]]")
	processedText = re.ReplaceAllLiteralString(processedText, " ")

	return TokenizedDocument{url: document.url, words: strings.Split(processedText, " ")}
}

/*
* @documents The list of documents to process
* Returns an array of processed documents
 */
func processDocuments(documents []*Document) []TokenizedDocument {

	var processors []TextProcessor
	var tokenizedDocs []TokenizedDocument

	downCaseProcessor := DownCaseProcessor{}
	accentProcessor := AccentProcessor{}
	processors = append(processors, downCaseProcessor, accentProcessor)

	for _, doc := range documents {
		tokenizedDocs = append(tokenizedDocs, analyse(*doc, processors))
	}

	return tokenizedDocs
}
