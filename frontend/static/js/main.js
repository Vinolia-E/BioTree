import { displayFileInfo, validFileType } from './file.js';

document.addEventListener("DOMContentLoaded", () => {
    const fileInput = document.getElementById("document-upload");
    const process = document.getElementById("processBtn");

    fileInput.addEventListener("change", (event) => {
        const file = event.target.files[0];

        if (!validFileType(file)) {
            alert("Invalid file type. Please upload a PDF, DOCX, or TXT file.");
            return;
        }

        if (file) {
            displayFileInfo(file);
        }


        process.disabled = false;
    });
});