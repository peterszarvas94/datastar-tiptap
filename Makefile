AIR ?= air
NPX ?= npx

.PHONY: dev css seed

dev:
	$(AIR)

css:
	$(NPX) tailwindcss -i ./css/input.css -o ./static/generated.css --watch

seed:
	mkdir -p data
	sqlite3 data/content.db < seed.sql
