package main

import (
	"encoding/gob"
	"os"
	"regexp"
	"strings"
	"text/scanner"
)

/*
* Structs
 */

/*
 * AST Grammar
 *
 *
 */

// Ast represents the parse tree
type Ast struct {
	root Node
}

// Node is a interface for all the nodes composing the ast
type Node interface {
	eval(index *Index) []int
}

// BinOpNode represents nodes with two operands
type BinOpNode struct {
	opType string
	left   *Node
	right  *Node
}

// UnOpNode represents nodes with one operand
type UnOpNode struct {
	opType string
	left   *Node
}

// WordNode represents nodes with a value
type WordNode struct {
	value string
}

func (nd BinOpNode) eval(index *Index) []int {
	var result []int
	leftEval := (*nd.left).eval(index)
	rightEval := (*nd.right).eval(index)

	if strings.Compare(nd.opType, "and") == 0 {
		result = intersection(leftEval, rightEval)
	} else {
		result = union(leftEval, rightEval)
	}

	return result
}

func (nd WordNode) eval(index *Index) []int {
	return index.wordToDids[nd.value]
}

func (nd UnOpNode) eval(index *Index) []int {
	var result []int
	urlToDid := index.urlToDid
	dids := (*nd.left).eval(index)
	didsMap := make(map[int]bool)

	for _, value := range dids {
		didsMap[value] = true
	}

	for _, v := range urlToDid {
		_, prs := didsMap[v]
		if prs {
			result = append(result, v)
		}
	}

	return result
}

func load(path string) Index {
	var urlMap map[string]int
	var didMap map[int]string
	var didsMap map[string][]int
	f, err := os.Open(path)
	defer f.Close()
	handleError(err, "Error while loading the index")

	dec := gob.NewDecoder(f)
	err = dec.Decode(&urlMap)
	err = dec.Decode(&didMap)
	err = dec.Decode(&didsMap)
	handleError(err, "Error while loading the index")
	return Index{urlToDid: urlMap, didToUrl: didMap, wordToDids: didsMap}
}

func makeAST(expression string) Node {
	var sc scanner.Scanner
	var root Node
	sc.Init(strings.NewReader(expression))
	root = checkExpRule(&sc)
	return root
}

func checkExpRule(sc *scanner.Scanner) Node {
	var tok string
	var r rune
	nd := checkTRule(sc)

	r = (*sc).Scan()
	if r == scanner.EOF {
		return nd
	}

	tok = (*sc).TokenText()
	if strings.Compare(tok, "and") == 0 ||
		strings.Compare(tok, "or") == 0 {
		var bin BinOpNode
		bin.opType = tok
		bin.left = &nd
		rightNode := checkExpRule(sc)
		bin.right = &rightNode

		if bin.right == nil {
			panic("Error while parsing the AST -- Exp rule")
		}

		return bin
	}

	return nil
}

func checkTRule(sc *scanner.Scanner) Node {
	var tok string
	re := regexp.MustCompile("[[:word:]]")
	_ = (*sc).Scan()
	tok = (*sc).TokenText()

	if re.MatchString(tok) {
		nd := WordNode{value: tok}
		return nd
	}

	panic("Syntax Error -- T rule")
}

func search(word string) []string {
	var ast Ast

	word = strings.ToLower(word)
	index := load(pathToIndexSave)
	ast.root = makeAST(word)
	docs := astSearch(ast, &index)

	return docs
}

func astSearch(ast Ast, index *Index) []string {
	var result []string
	dids := ast.root.eval(index)

	for _, value := range dids {
		result = append(result, index.didToUrl[value])
	}

	return result
}
