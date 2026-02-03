.PHONY: dev build

BASE_PATH ?= /demo/datastar-tiptap

dev:
	air & npx tailwindcss -i ./css/input.css -o ./static/generated.css --watch

build:
	npm install
	npx tailwindcss -i ./css/input.css -o ./static/generated.css
	BASE_PATH=$(BASE_PATH) go build -o ./bin/app ./src
