import { displayFileInfo, validFileType } from './file.js';
import { process  } from './process.js';

document.addEventListener("DOMContentLoaded", () => {
    const fileInput = document.getElementById("document-upload");
    const processBtn = document.getElementById("processBtn");
    let doc

    fileInput.addEventListener("change", (event) => {
        const file = event.target.files[0];

        if (!validFileType(file)) {
            alert("Invalid file type. Please upload a PDF, DOCX, or TXT file.");
            return;
        }

        if (file) {
            displayFileInfo(file);
            doc = file;
        } else {
            alert("No file selected. Please choose a file to upload.");
            return;
        }


        processBtn.disabled = false;
    });

    processBtn.addEventListener("click", (e) => {
        e.preventDefault();
        let formdata = new FormData();
        formdata.append("document", doc);
        let data = process(formdata);
        console.log(data);
    })
});
