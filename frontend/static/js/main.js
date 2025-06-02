import { displayFileInfo, validFileType } from './file.js';
import { process } from './process.js';
import { showNotification } from "./notification";

document.addEventListener("DOMContentLoaded", () => {
    const fileInput = document.getElementById("document-upload");
    const processBtn = document.getElementById("processBtn");
    let doc;
    let processedDataFile = "";

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
            
            if (data && data.units && data.dataFile) {
                processedDataFile = data.dataFile; // Store the data file name
                const container = document.getElementById("controls-container"); 
                container.style.display = "block";
                container.innerHTML = ""; 

                const select = document.createElement("select");
                select.id = "contentinfo";

                data.units.forEach(unit => {
                    const option = document.createElement("option");
                    option.value = unit;
                    option.textContent = unit;
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
                            <input type="number" id="chart-width" min="300" max="1200" value="800">
                        </div>
                        <div class="control-row">
                            <label for="chart-height">Height:</label>
                            <input type="number" id="chart-height" min="200" max="800" value="400">
                        </div>
                        <div class="control-row">
                            <label for="chart-type">Chart Type:</label>
                            <select id="chart-type">
                                <option value="line">Line Chart</option>
                                <option value="bar">Bar Chart</option>
                                <option value="pie">Pie Chart</option>
                            </select>
                        </div>
                        <button id="update-chart-btn">Show Chart</button>
                    </div>
                    <div class="control-group">
                        <h3>Actions</h3>
                        <button id="download-svg-btn" disabled>Download SVG</button>
                        <button id="reset-chart-btn">Reset Chart</button>
                    </div>
                    <div id="chart-container" style="margin-top: 20px;"></div>
                `);

                // Add event listener for the Show Chart button
                const updateChartBtn = document.getElementById("update-chart-btn");
                updateChartBtn.addEventListener("click", generateChart);

                // Add event listener for Download SVG button
                const downloadSvgBtn = document.getElementById("download-svg-btn");
                downloadSvgBtn.addEventListener("click", downloadSVG);

                // Add event listener for Reset Chart button
                const resetChartBtn = document.getElementById("reset-chart-btn");
                resetChartBtn.addEventListener("click", resetChart);
            }
        }).catch((error) => {
            console.error("Unexpected error:", error);
        });
    });

    // Function to generate chart
    async function generateChart() {
        const selectedUnit = document.getElementById("contentinfo").value;
        const chartTitle = document.getElementById("chart-title").value;
        const xLabel = document.getElementById("chart-x-label").value;
        const yLabel = document.getElementById("chart-y-label").value;
        const width = parseInt(document.getElementById("chart-width").value) || 800;
        const height = parseInt(document.getElementById("chart-height").value) || 400;
        const chartType = document.getElementById("chart-type").value;

        if (!selectedUnit) {
            showNotification("Please select a unit first.");
            return;
        }

        // Show loading state
        const updateChartBtn = document.getElementById("update-chart-btn");
        const originalText = updateChartBtn.textContent;
        updateChartBtn.textContent = "Generating...";
        updateChartBtn.disabled = true;

        const chartRequest = {
            data_file: processedDataFile,
            unit: selectedUnit,
            chart_type: chartType,
            title: chartTitle,
            x_label: xLabel,
            y_label: yLabel,
            width: width,
            height: height
        };

        try {
            const response = await fetch("/generate-chart", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    "Accept": "application/json"
                },
                body: JSON.stringify(chartRequest)
            });

            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            const result = await response.json();

            if (result.status === "error") {
                throw new Error(result.message || "Failed to generate chart");
            }

            // Display the chart
            const chartContainer = document.getElementById("chart-container");
            chartContainer.innerHTML = result.svg;

            // Enable download button
            document.getElementById("download-svg-btn").disabled = false;
            
            // Store SVG for download
            window.currentSVG = result.svg;

        } catch (error) {
            console.error("Error generating chart:", error);
            showNotification("Error generating chart: " + error.message);
        } finally {
            // Reset button state
            updateChartBtn.textContent = originalText;
            updateChartBtn.disabled = false;
        }
    }

    // Function to download SVG
    function downloadSVG() {
        if (!window.currentSVG) {
            showNotification("No chart available to download.");
            return;
        }

        const blob = new Blob([window.currentSVG], { type: "image/svg+xml" });
        const url = URL.createObjectURL(blob);
        const link = document.createElement("a");
        link.href = url;
        link.download = "chart.svg";
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
        URL.revokeObjectURL(url);
    }

    // Function to reset chart
    function resetChart() {
        const chartContainer = document.getElementById("chart-container");
        chartContainer.innerHTML = "";
        document.getElementById("download-svg-btn").disabled = true;
        window.currentSVG = null;
        
        // Reset form fields
        document.getElementById("chart-title").value = "";
        document.getElementById("chart-x-label").value = "";
        document.getElementById("chart-y-label").value = "";
        document.getElementById("chart-width").value = "800";
        document.getElementById("chart-height").value = "400";
        document.getElementById("chart-type").value = "line";
    }
});
