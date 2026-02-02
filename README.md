# Datastar with Tiptap

Minimal setup, intended for copy and pasting to other projects.

## Frontend

### CSS

Tailwind CSS and daisyUI

### JS

It uses importmaps with pre-bundled static files.

- Datastar source copied from a CDN, already ESM
- Tiptap editor pre-built into ESM
- JS-beautify pre-built into ESM

Bundle commands:

- `npm run bundle:tiptap` -> `static/vendor/tiptap.js`
- `npm run bundle:js-beautify` -> `static/vendor/js-beautify.js`

Other than that pre-bundling, no build step for JS, or runtime dependencies.
