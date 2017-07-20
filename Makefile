SRC= $(wildcard src/*.go)
BIN=myGoogle

all: build

build:
	go build -o ${BIN} ${SRC}

debug:
	go build -gcflags "-N -l" ${SRC}

install:
	go get golang.org/x/text/transform 
	go get golang.org/x/text/unicode/norm

clean:
	${RM} ${BIN}
	${RM} *.idx

.PHONY: clean install
