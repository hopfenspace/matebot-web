# Go compiler
GO = go

# Directory containing Go source files
FILE_DIR = .

# Directory containing compiled executables
OUT_DIR = bin

# File suffix for Go source files
FILE_SUFFIX = go

all: clean build

.PHONY: clena
clean:
	rm -rf ${OUT_DIR}/

.PHONY: build
build:
	${GO} build -o ${OUT_DIR}/ ${FILE_DIR}/...

.PHONY: install
install:

