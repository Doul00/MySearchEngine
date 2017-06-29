SRC= src/main.go src/document.go src/indexer.go src/searcher.go src/sliceUtils.go

all: build

build:
	go build ${SRC}

debug:
	go build -gcflags "-N -l" ${SRC}

install:
	go get golang.org/x/text/transform 
	go get golang.org/x/text/unicode/norm

clean:
	rm -rf bin/*

.PHONY: clean install
