import { displayFileInfo, validFileType } from './file.js';
import { process } from './process.js';
import { showNotification } from "./notification.js";
import { fetchUserDocuments, displayUserDocuments } from './documents.js';

document.addEventListener("DOMContentLoaded", () => {
    const fileInput = document.getElementById("document-upload");
    const processBtn = document.getElementById("processBtn");
    let doc;
    let processedDataFile = "";

    // Initialize user documents section
    initUserDocuments();

    fileInput.addEventListener("change", (event) => {
        const file = event.target.files[0];

        if (!validFileType(file)) {
            showNotification("Invalid file type. Please upload a PDF, DOCX, or TXT file.");
            return;
        }

        if (file) {
            displayFileInfo(file);
            doc = file;
        } else {
            showNotification("No file selected. Please choose a file to upload.");
            return;
        }

        processBtn.disabled = false;
    });

    processBtn.addEventListener("click", (e) => {
        e.preventDefault();
        
        if (!doc) {
            showNotification("No file selected. Please choose a file first.");
            return;
        }

        // Show loading state
        processBtn.textContent = "Processing...";
        processBtn.disabled = true;

        const formdata = new FormData();
        formdata.append("document", doc);

        console.log("Starting document processing...");
        
        process(formdata).then((data) => {
            console.log("Process response:", data);
            
            if (data) {
                // Reset the file input
                fileInput.value = "";
                doc = null;
                
                // Update the user documents list
                initUserDocuments();
                
                // Show success message
                showNotification("Document processed successfully!");
            }
            
            // Reset button state
            processBtn.textContent = "Process Document";
            processBtn.disabled = true;
        });
    });

    // Function to initialize user documents section
    async function initUserDocuments() {
        const documents = await fetchUserDocuments();
        displayUserDocuments(documents);
    }
});
