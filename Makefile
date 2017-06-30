SRC= $(wildcard src/*.go)
BIN=main

all: build

build:
	go build -o ${BIN} ${SRC}

debug:
	go build -gcflags "-N -l" ${SRC}

install:
	go get golang.org/x/text/transform 
	go get golang.org/x/text/unicode/norm

clean:
	rm ${BIN}

.PHONY: clean install
