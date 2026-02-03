# Datastar with Tiptap

Minimal setup, intended for copy and pasting to other projects.

## Tooling

This repo uses `mise` to pin tool versions. If you don't use `mise`, make sure you have these installed:

- Go (1.25.x)
- Node.js (25.x)

Node modules are only used for Tailwind tooling and for bundling/copying vendor packages into `static/vendor`.

The server loads environment variables from a local `.env` file (optional).

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

- `make dev` -> runs `air` plus the Tailwind CSS watcher

## Environment

- `BASE_PATH` -> optional URL prefix (example: `/demo/datastar-tiptap`)

## Go Layout

Go source files live in `src`.
