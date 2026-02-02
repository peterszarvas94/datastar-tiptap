import "datastar";
import { Editor, StarterKit } from "tiptap-bundle";
import jsBeautify from "js-beautify";

// Initialize editor
const editor = new Editor({
  element: document.querySelector("#editor"),
  extensions: [StarterKit.configure({ underline: true })],
  content: "", // backend fills this
  onUpdate: () => {
    window.dispatchEvent(new CustomEvent("editorupdate"));
  },
  onSelectionUpdate: () => {
    window.dispatchEvent(new CustomEvent("editorupdate"));
  },
});

window.editor = editor;
window.beautify = jsBeautify.html_beautify;
