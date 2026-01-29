AIR ?= air
NPX ?= npx

.PHONY: dev css

dev:
	$(AIR)

css:
	$(NPX) tailwindcss -i ./css/input.css -o ./static/generated.css
