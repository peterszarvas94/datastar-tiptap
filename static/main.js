import "datastar";
import { Editor, StarterKit } from "tiptap-bundle";
import htmlBeautify from "js-beautify";

const API_URL = "";

// Initialize editor
const editor = new Editor({
  element: document.querySelector("#editor"),
  extensions: [StarterKit],
  content: "", // backend fills this
  onUpdate: () => {
    updateToolbarState();
  },
  onSelectionUpdate: () => {
    updateToolbarState();
  },
});

// Update active states
function updateToolbarState() {
  toggleToolbarButton("#btn-bold", editor.isActive("bold"));
  toggleToolbarButton("#btn-italic", editor.isActive("italic"));
  toggleToolbarButton("#btn-strike", editor.isActive("strike"));
  const isParagraph =
    editor.isActive("paragraph") &&
    !editor.isActive("bulletList") &&
    !editor.isActive("orderedList") &&
    !editor.isActive("blockquote") &&
    !editor.isActive("codeBlock");
  toggleToolbarButton("#btn-paragraph", isParagraph);
  toggleToolbarButton("#btn-h1", editor.isActive("heading", { level: 1 }));
  toggleToolbarButton("#btn-h2", editor.isActive("heading", { level: 2 }));
  toggleToolbarButton("#btn-h3", editor.isActive("heading", { level: 3 }));
  toggleToolbarButton("#btn-bullet", editor.isActive("bulletList"));
  toggleToolbarButton("#btn-ordered", editor.isActive("orderedList"));
  toggleToolbarButton("#btn-blockquote", editor.isActive("blockquote"));
  toggleToolbarButton("#btn-code", editor.isActive("codeBlock"));
}

function toggleToolbarButton(selector, isActive) {
  const button = document.querySelector(selector);
  if (!button) return;
  button.classList.toggle("btn-secondary", isActive);
}

window.editor = editor;
window.beautify = htmlBeautify.html_beautify;
