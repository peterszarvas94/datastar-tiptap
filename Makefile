.PHONY: dev css

dev:
	go run ./src

css:
	npx tailwindcss -i ./css/input.css -o ./static/generated.css --watch
