package main

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"io"

	"encoding/json"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

/*
* Structs
 */

// Document contains the information regarding a document
type Document struct {
	url, title, body string
}

// TokenizedDocument contains the document and all its words
type TokenizedDocument struct {
	Title, Body map[string]int
	Url         string
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

func createDocument(rawDoc string) Document {
	var doc Document
	json.Unmarshal([]byte(rawDoc), &doc)
	return doc
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
func applyProcessors(document Document, processors []TextProcessor) TokenizedDocument {
	processedTitle := document.title
	processedBody := document.body

	for _, p := range processors {
		p.process(&processedTitle)
		p.process(&processedBody)
	}

	re := regexp.MustCompile("[[:^word:]]")
	processedTitle = re.ReplaceAllLiteralString(processedTitle, " ")
	processedBody = re.ReplaceAllLiteralString(processedBody, " ")

	return TokenizedDocument{
		Title: createCounter(processedTitle),
		Body:  createCounter(processedBody),
		Url:   document.url}
}

func createCounter(text string) map[string]int {
	splitText := strings.Split(text, " ")
	// pre-allocation to half the size of the vocab (with duplicata).
	counter := make(map[string]int, int(len(splitText)/2))
	for _, word := range splitText {
		counter[word]++
	}

	return counter
}

/*
* @documents The list of documents to process
* Returns an array of processed documents
 */
func processDocuments(documents *[]Document) []TokenizedDocument {

	var processors []TextProcessor
	var tokenizedDocs []TokenizedDocument

	downCaseProcessor := DownCaseProcessor{}
	accentProcessor := AccentProcessor{}
	processors = append(processors, downCaseProcessor, accentProcessor)

	for _, doc := range *documents {
		tokenizedDocs = append(tokenizedDocs, applyProcessors(doc, processors))
	}

	return tokenizedDocs
}
