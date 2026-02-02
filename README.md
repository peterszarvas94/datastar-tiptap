# Datastar with Tiptap

Minimal setup, intended for copy and pasting to other projects.

## Tooling

This repo uses `mise` to pin tool versions. If you don't use `mise`, make sure you have these installed:

- Go (1.25.x)
- Node.js (25.x)
- Air (1.64.x)

Node modules are only used for Tailwind tooling and for bundling/copying vendor packages into `static/vendor`.

## Frontend

### CSS

Tailwind CSS and daisyUI

### JS

It uses importmaps with pre-bundled static files.

- Datastar source copied from a CDN, already ESM
- Tiptap editor pre-built into ESM
- JS-beautify pre-built into ESM
- entry point: `static/main.js`

Bundle commands:

- `npm run bundle:tiptap` -> `static/vendor/tiptap.js`
- `npm run bundle:js-beautify` -> `static/vendor/js-beautify.js`

Other than that pre-bundling, no build step for JS, or runtime dependencies.

## Make Commands

- `make dev` -> runs the Go server with live reload (air)
- `make css` -> watches and rebuilds Tailwind CSS
- `make seed` -> seeds SQLite content from `seed.sql`
