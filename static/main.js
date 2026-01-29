import { Editor } from "@tiptap/core";
import StarterKit from "@tiptap/starter-kit";
import htmlBeautify from "js-beautify";

const API_URL = "";

// API functions
async function getContent() {
  try {
    const response = await fetch(`${API_URL}/api/content`);
    const data = await response.json();
    return data.html || "";
  } catch (error) {
    console.error("Failed to load content:", error);
    return "";
  }
}

async function saveContent(html) {
  try {
    const response = await fetch(`${API_URL}/api/content`, {
      method: "PATCH",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ html }),
    });
    if (!response.ok) {
      throw new Error("Failed to save content");
    }
    return true;
  } catch (error) {
    console.error("Failed to save content:", error);
    return false;
  }
}

// Initialize editor
const editor = new Editor({
  element: document.querySelector("#editor"),
  extensions: [StarterKit],
  content: "",
  onUpdate: () => {
    updateToolbarState();
  },
  onSelectionUpdate: () => {
    updateToolbarState();
  },
});

// Load initial content
(async () => {
  const content = await getContent();
  editor.commands.setContent(content);
  updateToolbarState();
  await updateDbPreview();
})();

// Button click handlers
const getBtn = (id) => {
  const btn = document.querySelector(id);
  if (!btn) throw new Error(`Button ${id} not found`);
  return btn;
};

getBtn("#btn-bold").addEventListener("click", () => {
  editor.chain().focus().toggleBold().run();
});
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

// Save button handler
getBtn("#btn-save").addEventListener("click", async () => {
  let html = editor.getHTML();

  // Remove trailing empty paragraphs
  html = html.replace(/(<p><\/p>\s*)+$/, "");

  const formatted = htmlBeautify.html_beautify(html, { indent_size: 2 });
  const success = await saveContent(formatted);
  if (success) {
    await updateDbPreview();
  } else {
    console.error("Save failed");
  }
});

// Copy button handler
getBtn("#btn-copy").addEventListener("click", async () => {
  const content = await getContent();
  try {
    await navigator.clipboard.writeText(content);
    // Optional: Show feedback (could add a toast notification later)
    const btn = getBtn("#btn-copy");
    const originalTitle = btn.title;
    btn.title = "Copied!";
    setTimeout(() => {
      btn.title = originalTitle;
    }, 2000);
  } catch (err) {
    console.error("Failed to copy:", err);
  }
});

// Update database preview
async function updateDbPreview() {
  const preview = document.querySelector("#db-preview");
  if (!preview) return;
  const content = await getContent();
  preview.textContent = content;
  updateRenderedPreview(content);
}

// Update rendered preview
function updateRenderedPreview(content) {
  const renderedPreview = document.querySelector("#rendered-preview");
  if (!renderedPreview) return;
  renderedPreview.innerHTML = content;
}

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
