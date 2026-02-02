// Entry point for bundling tiptap to static/vendor/tiptap.js.
// Keeps the esbuild command simple and makes it easy to add/remove
// exported extensions in one place without touching the bundle script.
//
// Instructions:
// - Add new exports here (e.g. extensions) to include them in the bundle.
// - Rebuild with: npm run bundle:tiptap (or npm run bundle:vendor).
export { Editor } from "@tiptap/core";
export { StarterKit } from "@tiptap/starter-kit";
export { Strike } from "@tiptap/extension-strike";
export { Underline } from "@tiptap/extension-underline";
