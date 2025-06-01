import { dom } from './dom.js';
import { state } from './state.js';

export function initEventListeners() {
  ["dragenter", "dragover", "dragleave", "drop"].forEach(eventName => {
    dom.dropArea.addEventListener(eventName, preventDefaults, false);
  });

  ["dragenter", "dragover"].forEach(eventName => {
    dom.dropArea.addEventListener(eventName, () => dom.dropArea.classList.add("highlight"), false);
  });

  ["dragleave", "drop"].forEach(eventName => {
    dom.dropArea.addEventListener(eventName, () => dom.dropArea.classList.remove("highlight"), false);
  });

  dom.dropArea.addEventListener("drop", handleDrop, false);
}

function preventDefaults(e) {
  e.preventDefault();
  e.stopPropagation();
}

function handleDrop(e) {
  const dt = e.dataTransfer;
  const files = dt.files;
  if (files.length) {
    const file = files[0];
    if (state.fileTypes.includes(file.type)) {
      dom.filePreview.innerHTML = `<p>Selected file: ${file.name}</p>`;
      dom.processBtn.disabled = false;
      state.currentDataFile = file;
      dom.visualizationSection.style.display = "block";
    } else {
      dom.filePreview.innerHTML = "<p>Unsupported file type.</p>";
    }
  }
}