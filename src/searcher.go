package main

import (
	"regexp"
	"strings"
)

/*
* Structs
 */

/*
 * AST Grammar
 * Exp := T | T and Exp | T or Exp
 * T := word | '(' Exp ')'| not Exp
 */

// Ast represents the parse tree
type Ast struct {
	root Node
}

// Node is a interface for all the nodes composing the ast
type Node interface {
	eval(index *Index) []string
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

/*
* @index the index
* Evaluates the expression of the node, based on its values and its sons values
* an 'and' node will intersect the results, an 'or' node will join them
 */
func (nd BinOpNode) eval(index *Index) []string {
	var urls []string
	leftEval := (*nd.left).eval(index)
	rightEval := (*nd.right).eval(index)

	if strings.Compare(nd.opType, "and") == 0 {
		urls = intersection(leftEval, rightEval)
	} else {
		urls = union(leftEval, rightEval)
	}

	return urls
}

/*
* @index the index
* Returns the value of the word node
 */
func (nd WordNode) eval(index *Index) []string {
	return index.Posting[nd.value]
}

/*
* @index the index
* Negates the results of the node's son
 */
func (nd UnOpNode) eval(index *Index) []string {
	var invertedResults []string

	rawResults := make(map[string]bool)
	for _, url := range (*nd.left).eval(index) {
		rawResults[url] = true
	}

	for url := range index.Posting {
		_, prs := rawResults[url]
		if !prs {
			invertedResults = append(invertedResults, url)
		}
	}

	return invertedResults
}

/*
* @expression the user's query
* Returns a Ast built from the query
 */
func makeAST(expression string) Node {
	var root Node
	var i = 0
	root = checkExpRule(strings.Split(expression, " "), &i)
	return root
}

/*
* @tokens the list of tokens
* @i the reading position.
* I did not use Go's Scanner because it does not have an 'unread' operation to go backwards
* Returns an Exp node
 */
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

/*
* @tokens the list of tokens
* @i the reading position.
* Returns a T node, containing either a terminal or other expressions
 */
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

/*
* @word the word to search for
* @index the index containing the documents information
* Builds the Ast from the query and returns the results by evaluating it
 */
func search(query string, index *Index) []TokenizedDocument {
	var ast Ast

	processors := []TextProcessor{DownCaseProcessor{}, AccentProcessor{}}
	processedQuery := processText(query, processors)
	ast.root = makeAST(processedQuery)
	docs := astSearch(ast, index)

	return docs
}

/*
* @ast the ast to evaluate
* @index the index
* Evaluates the ast and returns the results
 */
func astSearch(ast Ast, index *Index) []TokenizedDocument {
	var result []TokenizedDocument

	for _, url := range ast.root.eval(index) {
		result = append(result, index.Url2docs[url])
	}

	return result
}
