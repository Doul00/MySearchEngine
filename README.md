# MY SEARCH ENGINE
## SCIA - TMLN Class

This program is a little search engine written in Go.
Please install go before running it.
It works in two parts.

1. First, you need to create an index from your files. In order to do this, type:
./main --create-index path_to_your_folder index_filename

**path_to_your_folder** is the path to the folder where your text files are
**index_filename** is the name of the file containing the index. The data is serialized using Go's 'gob' serializer

2. Once the indexing is finished, To launch the search engine, type:
./main --search index_filename

**index_filename** is the name of the file containing the index.

The prompt will appear.

3. Type in your query following these guidelines. Do it **without the brackets**:

⋅⋅* To query a simple word, just type it
⋅⋅* You can use the "and" operator to query multiple words: 
For instance "Sylvain and Coca" (without the brackets)
⋅⋅* You can use the "or" operator:
For instance "sport or macdonalds" 
⋅⋅* You can use the "not" operator:
For instance "not coding"
