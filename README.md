# MY SEARCH ENGINE
## SCIA - TMLN Class

This program is a little search engine written in Go.
Please install Go before running it.
To download the external librairies used for the processors, type in:
- make install
If it does not work please type in your console:
	go get golang.org/x/text/transform 
	go get golang.org/x/text/unicode/norm

The project works in two parts

1. First, you need to create an index from your files. In order to do this, type:
./main -mode=create -path=path_to_your_folder -index=index_filename

**path_to_your_folder** is the path to the folder where your text files are

**index_filename** is the name of the file containing the index. The data is serialized using Go's 'gob' serializer

2. Once the indexing is finished, To launch the search engine, type:
./main -mode=search -index=index_filename

**index_filename** is the name of the file containing the index.

The prompt will appear.

3. Type in your query following these guidelines. Do it **without the brackets**:

- To query a simple word, just type it
- You can use the "and" operator to query multiple words: 
For instance "Sylvain and Coca" (without the brackets)
- You can use the "or" operator:
For instance "sport or macdonalds" 
- You can use the "not" operator:
For instance "not coding"
- You can use the parenthesis, but **always** put a whitespace between the parenthesis and the word.
There is no lexer so without whitespaces there will be an error!
For instance "not ( hell and heaven )"
