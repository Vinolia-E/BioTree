import { displayFileInfo, validFileType } from './file.js';
import { process } from './process.js';

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
        const formdata = new FormData();
        formdata.append("document", doc);

        process(formdata).then((data) => {
            if (data) {
                const container = document.getElementById("controls-container"); 
                container.style.display = "block";
                container.innerHTML = ""; 

                const select = document.createElement("select");
                select.id = "contentinfo";

                data.forEach(ele => {
                    const option = document.createElement("option");
                    option.value = ele;
                    option.textContent = ele;
                    select.appendChild(option);
                });

                container.appendChild(select);

                container.insertAdjacentHTML("beforeend", `
                    <div class="control-group">
                        <h3>Chart Parameters</h3>
                        <div class="control-row">
                            <label for="chart-title">Title:</label>
                            <input type="text" id="chart-title" >
                        </div>
                        <div class="control-row">
                            <label for="chart-x-label">X-Axis Label:</label>
                            <input type="text" id="chart-x-label" >
                        </div>
                        <div class="control-row">
                            <label for="chart-y-label">Y-Axis Label:</label>
                            <input type="text" id="chart-y-label" >
                        </div>
                        <div class="control-row">
                            <label for="chart-width">Width:</label>
                            <input type="number" id="chart-width"  min="300" max="1200">
                        </div>
                        <div class="control-row">
                            <label for="chart-height">Height:</label>
                            <input type="number" id="chart-height"  min="200" max="800">
                        </div>
                        <button id="update-chart-btn">show Chart</button>
                    </div>
                    <div class="control-group">
                        <h3>Actions</h3>
                        <button id="download-svg-btn">Download SVG</button>
                        <button id="reset-chart-btn">Reset Chart</button>
                    </div>
                `);
            }

        }).catch((error) => {
            console.error("Unexpected error:", error);
        });
    });

});
