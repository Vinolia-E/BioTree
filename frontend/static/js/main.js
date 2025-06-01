import { displayFileInfo, validFileType } from './dom.js';

document.addEventListener("DOMContentLoaded", () => {
    const fileInput = document.getElementById("document-upload");

    fileInput.addEventListener("change", (event) => {
        const file = event.target.files[0];

        if (!validFileType(file)) {
            alert("Invalid file type. Please upload a PDF, DOCX, or TXT file.");
            return;
        }
        
        if (file) {
            displayFileInfo(file);
        }
    });
});