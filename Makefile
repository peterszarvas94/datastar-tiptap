.PHONY: dev css

dev:
	air

css:
	npx tailwindcss -i ./css/input.css -o ./static/generated.css --watch
