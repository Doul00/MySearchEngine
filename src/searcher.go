package main

import (
	"encoding/gob"
	"os"
	"regexp"
	"strings"
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
	return generationSearch((*index).Generations, nd.value)
}

func (nd UnOpNode) eval(index *Index) []int {
	var result []int
	urlToDid := index.UrlToDid
	dids := (*nd.left).eval(index)
	didsMap := make(map[int]bool)

	for _, value := range dids {
		didsMap[value] = true
	}

	for _, v := range urlToDid {
		_, prs := didsMap[v]
		if !prs {
			result = append(result, v)
		}
	}

	return result
}

func load(path string) Index {
	var index Index

	f, err := os.Open(path)
	handleError(err, "Error while loading the index")

	dec := gob.NewDecoder(f)
	err = dec.Decode(&index)
	handleError(err, "Error while loading the index")
	f.Close()
	return index
}

func generationSearch(generations []Generation, word string) []int {
	for i := len(generations) - 1; i >= 0; i-- {
		currMap := generations[i].WordsToDid
		val, prs := currMap[word]

		if prs {
			return val
		}
	}

	return []int{}
}

func makeAST(expression string) Node {
	var root Node
	var i = 0
	root = checkExpRule(strings.Split(expression, " "), &i)
	return root
}

func checkExpRule(tokens []string, i *int) Node {
	var tok string
	nd := checkTRule(tokens, i)
	if *i == len(tokens) {
		return nd
	}
	tok = tokens[*i]

	if strings.Compare(tok, "and") == 0 ||
		strings.Compare(tok, "or") == 0 {
		var bin BinOpNode
		bin.opType = tok
		bin.left = &nd
		(*i)++
		rightNode := checkExpRule(tokens, i)
		bin.right = &rightNode

		if bin.right == nil {
			panic("Error while parsing the AST -- Exp rule")
		}
		return bin
	}

	return nd
}

func checkTRule(tokens []string, i *int) Node {
	var tok string
	re := regexp.MustCompile("[[:word:]]")

	tok = tokens[*i]

	if re.MatchString(tok) {
		if strings.Compare(tok, "not") == 0 {
			nd := UnOpNode{opType: tok}
			(*i)++
			leftSon := checkExpRule(tokens, i)
			nd.left = &leftSon
			return nd
		} else {
			(*i)++
			nd := WordNode{value: tok}
			return nd
		}
	} else {
		if strings.Compare(tok, "(") == 0 {
			(*i)++
			nd := checkExpRule(tokens, i)
			tok = tokens[*i]
			if strings.Compare(tok, ")") == 0 {
				return nd
			}
		}
	}

	panic("Syntax Error -- T rule")
}

func search(word string, index Index) []string {
	var ast Ast

	word = strings.ToLower(word)
	ast.root = makeAST(word)
	docs := astSearch(ast, &index)

	return docs
}

func astSearch(ast Ast, index *Index) []string {
	var result []string
	dids := ast.root.eval(index)

	for _, value := range dids {
		result = append(result, index.DidToUrl[value])
	}

	return result
}
