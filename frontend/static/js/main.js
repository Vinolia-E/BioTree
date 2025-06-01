import { displayFileInfo } from './dom.js';

document.addEventListener("DOMContentLoaded", () => {
    const fileInput = document.getElementById("document-upload");

    fileInput.addEventListener("change", (event) => {
        const file = event.target.files[0];
        if (file) {
            displayFileInfo(file);
        }
    });
});