import "datastar";
import { Editor, StarterKit } from "tiptap-bundle";
import htmlBeautify from "js-beautify";

const API_URL = "";

// API functions
async function getContent() {
  try {
    // this returns html
    const response = await fetch(`${API_URL}/content`);
    const html = await response.text();
    return html || "";
  } catch (error) {
    console.error("Failed to load content:", error);
    return "";
  }
}

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

// Button click handlers
const getBtn = (id) => {
  const btn = document.querySelector(id);
  if (!btn) throw new Error(`Button ${id} not found`);
  return btn;
};

// getBtn("#btn-bold").addEventListener("click", () => {
//   editor.chain().focus().toggleBold().run();
// });
getBtn("#btn-italic").addEventListener("click", () => {
  editor.chain().focus().toggleItalic().run();
});
getBtn("#btn-strike").addEventListener("click", () => {
  editor.chain().focus().toggleStrike().run();
});
getBtn("#btn-paragraph").addEventListener("click", () => {
  editor.chain().focus().setParagraph().run();
});
getBtn("#btn-h1").addEventListener("click", () => {
  editor.chain().focus().toggleHeading({ level: 1 }).run();
});
getBtn("#btn-h2").addEventListener("click", () => {
  editor.chain().focus().toggleHeading({ level: 2 }).run();
});
getBtn("#btn-h3").addEventListener("click", () => {
  editor.chain().focus().toggleHeading({ level: 3 }).run();
});
getBtn("#btn-bullet").addEventListener("click", () => {
  editor.chain().focus().toggleBulletList().run();
});
getBtn("#btn-ordered").addEventListener("click", () => {
  editor.chain().focus().toggleOrderedList().run();
});
getBtn("#btn-blockquote").addEventListener("click", () => {
  editor.chain().focus().toggleBlockquote().run();
});
getBtn("#btn-code").addEventListener("click", () => {
  editor.chain().focus().toggleCodeBlock().run();
});
getBtn("#btn-undo").addEventListener("click", () => {
  editor.chain().focus().undo().run();
});
getBtn("#btn-redo").addEventListener("click", () => {
  editor.chain().focus().redo().run();
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
