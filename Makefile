.PHONY: dev build

dev:
	air & npx tailwindcss -i ./css/input.css -o ./static/generated.css --watch

build:
	npm install
	npx tailwindcss -i ./css/input.css -o ./static/generated.css
	go build -o ./bin/app ./src
