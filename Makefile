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
	useradd -U -r matebot-web
	cp ${OUT_DIR}/matebot-web /usr/local/bin/matebot-web
	cp matebot-web.service /usr/lib/systemd/system/matebot-web.service
	if [ -L /etc/systemd/system/multi-user.target.wants/matebot-web.service ] ; then \
		if [ -e /etc/systemd/system/multi-user.target.wants/matebot-web.service ]; then \
				echo "Service file is already linked properly"; \
		else \
			rm /etc/systemd/system/multi-user.target.wants/matebot-web.service; \
			ln -s /usr/lib/systemd/system/matebot-web.service /etc/systemd/system/multi-user.target.wants/; \
		fi \
	else \
		ln -s /usr/lib/systemd/system/matebot-web.service /etc/systemd/system/multi-user.target.wants/; \
	fi
	systemctl daemon-reload
	systemctl enable matebot-web
	cp -r templates/ /var/lib/matebot-web/
	cp -r static/ /var/lib/matebot-web/
	cp example.config.toml /etc/matebot-web/
