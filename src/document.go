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
	Url   string `json:"url"`
	Title string `json:"title"`
	Body  string `json:"body"`
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
func processDocument(document Document, processors []TextProcessor) TokenizedDocument {
	processedTitle := processText(document.Title, processors)
	processedBody := processText(document.Body, processors)

	return TokenizedDocument{
		Title: createCounter(processedTitle),
		Body:  createCounter(processedBody),
		Url:   document.Url}
}

func processText(text string, processors []TextProcessor) string {
	for _, p := range processors {
		p.process(&text)
	}

	re := regexp.MustCompile("[[:^word:]]")
	return re.ReplaceAllLiteralString(text, " ")
}

func createCounter(text string) map[string]int {
	splitText := strings.Split(text, " ")
	counter := make(map[string]int)
	for _, word := range splitText {
		if len(word) != 0 && word != " " {
			counter[word]++
		}
	}

	return counter
}

/*
* @documents The list of documents to process
* Returns an array of processed documents
 */
func processDocuments(documents *[]Document) []TokenizedDocument {
	var tokenizedDocs []TokenizedDocument

	processors := []TextProcessor{DownCaseProcessor{}, AccentProcessor{}}

	for _, doc := range *documents {
		tokenizedDocs = append(tokenizedDocs, processDocument(doc, processors))
	}

	return tokenizedDocs
}
