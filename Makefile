.PHONY: dev

dev:
	air & npx tailwindcss -i ./css/input.css -o ./static/generated.css --watch
