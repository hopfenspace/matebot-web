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
	npm init -y
	npm install babel-cli@6 babel-preset-react-app@3
	npx babel js_src --out-dir static/js --presets react-app/prod --minified --ignore=m_sdk.js
	cp js_src/lib/m_sdk.js static/js/lib/

.PHONY: install
install:

