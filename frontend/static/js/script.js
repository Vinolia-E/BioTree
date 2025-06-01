/**
 * BioTree - SVG Chart Generation and Visualization
 * Enhanced JavaScript for SVG handler integration
 */

document.addEventListener("DOMContentLoaded", () => {
  // DOM Elements
  const dropArea = document.getElementById("dropArea");
  const fileInput = document.getElementById("document-upload");
  const filePreview = document.getElementById("filePreview");
  const chartPreview = document.getElementById("chartPreview");
  const processBtn = document.getElementById("processBtn");
  const chartOptions = document.querySelectorAll('input[name="chart-type"]');
  const visualizationSection = document.querySelector(".visualization-section");

  // State management
  let currentSvg = null;
  let currentDataFile = null;
  let chartParameters = {
    title: "Data Visualization",
    x_label: "Categories",
    y_label: "Values",
    width: 800,
    height: 400,
  };

  // Accepted files
  const fileTypes = [
    "application/pdf",
    "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
    "text/plain",
  ];

  // Initialize with visualization section hidden
  visualizationSection.style.display = "none";

  // Event Listeners for drag and drop
  ["dragenter", "dragover", "dragleave", "drop"].forEach((eventName) => {
    dropArea.addEventListener(eventName, preventDefaults, false);
  });

  function preventDefaults(e) {
    e.preventDefault();
    e.stopPropagation();
  }

  ["dragenter", "dragover"].forEach((eventName) => {
    dropArea.addEventListener(eventName, highlight, false);
  });

  ["dragleave", "drop"].forEach((eventName) => {
    dropArea.addEventListener(eventName, unhighlight, false);
  });

  function highlight() {
    dropArea.classList.add("highlight");
  }

  function unhighlight() {
    dropArea.classList.remove("highlight");
  }

  // Handle dropped files
  dropArea.addEventListener("drop", handleDrop, false);

  function handleDrop(e) {
    const dt = e.dataTransfer;
    const files = dt.files;

    if (files.length) {
      fileInput.files = files;
      handleFiles(files);
    }
  }

  // Handle file selection via input
  fileInput.addEventListener("change", function () {
    handleFiles(this.files);
  });

  // Process the selected files
  function handleFiles(files) {
    const file = files[0];

    if (validFileType(file)) {
      displayFileInfo(file);
      processBtn.disabled = false;
      visualizationSection.style.display = "block";
      updateChartPreview();
    } else {
      filePreview.innerHTML = `
        <p class="error">Error: ${file.name} is not a supported file type. 
        Please select a PDF, DOCX, or TXT file.</p>
      `;
      processBtn.disabled = true;
      visualizationSection.style.display = "none";
    }
  }

  // Check if file type is valid
  function validFileType(file) {
    return (
      fileTypes.includes(file.type) ||
      file.name.endsWith(".pdf") ||
      file.name.endsWith(".docx") ||
      file.name.endsWith(".txt")
    );
  }

  // Display file information
  function displayFileInfo(file) {
    // Format file size
    const size = formatFileSize(file.size);

    // Get file icon based on type
    const iconType = getFileIconType(file);

    filePreview.innerHTML = `
      <div class="file-metadata">
        <div class="file-icon ${iconType}"></div>
        <div class="file-info">
          <p><strong>File name:</strong> ${file.name}</p>
          <p><strong>File size:</strong> ${size}</p>
          <p><strong>Last modified:</strong> ${new Date(
            file.lastModified
          ).toLocaleDateString()}</p>
        </div>
      </div>
    `;
  }

  // Format file size
  function formatFileSize(bytes) {
    if (bytes < 1024) {
      return bytes + " bytes";
    } else if (bytes < 1048576) {
      return (bytes / 1024).toFixed(1) + " KB";
    } else {
      return (bytes / 1048576).toFixed(1) + " MB";
    }
  }

  // Get file icon type
  function getFileIconType(file) {
    if (file.type === "application/pdf" || file.name.endsWith(".pdf")) {
      return "pdf-icon";
    } else if (
      file.type ===
        "application/vnd.openxmlformats-officedocument.wordprocessingml.document" ||
      file.name.endsWith(".docx")
    ) {
      return "docx-icon";
    } else {
      return "txt-icon";
    }
  }

  // Update chart preview based on selected chart type
  function updateChartPreview() {
    const selectedChart = document.querySelector(
      'input[name="chart-type"]:checked'
    ).value;

    if (currentDataFile) {
      // If we have a data file, generate a new chart with the selected type
      generateChart(currentDataFile, selectedChart);
    } else {
      // Otherwise show a placeholder
      chartPreview.innerHTML = `
        <div class="chart-placeholder ${selectedChart}-placeholder">
          <p>Preview of ${selectedChart} chart</p>
          <p>This is where the ${selectedChart} chart will be rendered</p>
        </div>
      `;
    }
  }

  // Listen for chart type changes
  chartOptions.forEach((option) => {
    option.addEventListener("change", updateChartPreview);
  });

  // Process button click handler
  processBtn.addEventListener("click", async () => {
    const file = fileInput.files[0];
    if (!file) {
      showNotification("Please select a file first", "error");
      return;
    }

    const chartType = document.querySelector(
      'input[name="chart-type"]:checked'
    ).value;

    // Show loading state
    processBtn.disabled = true;
    processBtn.textContent = "Processing...";
    showLoadingIndicator(chartPreview);

    try {
      const result = await uploadAndProcessFile(file, chartType);

      // Store the data file name for later use
      currentDataFile = result.data_file;

      // Display the generated SVG chart
      if (result.svg) {
        displaySvgChart(result.svg);
        currentSvg = result.svg;

        // Add chart controls after successful generation
        addChartControls();
      }

      // Show success message
      showNotification(
        `File processed successfully! Chart type: ${result.chart_type}`,
        "success"
      );
    } catch (error) {
      console.error("Error:", error);
      showNotification("Error processing file. Please try again.", "error");
    } finally {
      // Reset button state
      processBtn.disabled = false;
      processBtn.textContent = "Process Document";
    }
  });

  // Upload and process file
  async function uploadAndProcessFile(file, chartType) {
    const formData = new FormData();
    formData.append("document", file);
    formData.append("chart_type", chartType);

    const response = await fetch("/api/process-and-generate", {
      method: "POST",
      body: formData,
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      throw new Error(errorData.message || "Network response was not ok");
    }

    return await response.json();
  }

  // Generate chart from existing data file
  async function generateChart(dataFile, chartType) {
    if (!dataFile) return;

    showLoadingIndicator(chartPreview);

    try {
      const response = await fetch("/api/generate-chart", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          data_file: dataFile,
          chart_type: chartType,
          title: chartParameters.title,
          x_label: chartParameters.x_label,
          y_label: chartParameters.y_label,
          width: chartParameters.width,
          height: chartParameters.height,
        }),
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(errorData.message || "Failed to generate chart");
      }

      const result = await response.json();

      if (result.svg) {
        displaySvgChart(result.svg);
        currentSvg = result.svg;
      }
    } catch (error) {
      console.error("Error generating chart:", error);
      showNotification("Error generating chart. Please try again.", "error");
    }
  }

  // Display SVG chart
  function displaySvgChart(svgContent) {
    chartPreview.innerHTML = svgContent;

    // Make SVG responsive
    const svgElement = chartPreview.querySelector("svg");
    if (svgElement) {
      svgElement.setAttribute("width", "100%");
      svgElement.setAttribute("height", "auto");
      svgElement.style.maxHeight = "500px";

      // Add interactivity to SVG elements
      enhanceSvgInteractivity(svgElement);
    }
  }

  // Add interactivity to SVG elements
  function enhanceSvgInteractivity(svgElement) {
    // Add tooltips to data points
    const dataElements = svgElement.querySelectorAll("circle, rect");
    dataElements.forEach((element) => {
      element.style.cursor = "pointer";

      // Create tooltip element
      const tooltip = document.createElement("div");
      tooltip.className = "svg-tooltip";
      tooltip.style.position = "absolute";
      tooltip.style.display = "none";
      tooltip.style.background = "rgba(0,0,0,0.8)";
      tooltip.style.color = "white";
      tooltip.style.padding = "5px";
      tooltip.style.borderRadius = "3px";
      tooltip.style.fontSize = "12px";
      tooltip.style.pointerEvents = "none";
      tooltip.style.zIndex = "1000";
      document.body.appendChild(tooltip);

      // Get data values from element attributes
      const value =
        element.getAttribute("data-value") ||
        (element.getAttribute("height")
          ? parseFloat(element.getAttribute("height")).toFixed(2)
          : "N/A");

      const category = element.getAttribute("data-category") || "Data point";

      // Show tooltip on hover
      element.addEventListener("mouseover", (e) => {
        tooltip.textContent = `${category}: ${value}`;
        tooltip.style.display = "block";
        tooltip.style.left = `${e.pageX + 10}px`;
        tooltip.style.top = `${e.pageY + 10}px`;

        // Highlight the element
        element.setAttribute(
          "data-original-fill",
          element.getAttribute("fill")
        );
        element.setAttribute("fill", "#ff7700");
      });

      // Move tooltip with cursor
      element.addEventListener("mousemove", (e) => {
        tooltip.style.left = `${e.pageX + 10}px`;
        tooltip.style.top = `${e.pageY + 10}px`;
      });

      // Hide tooltip when not hovering
      element.addEventListener("mouseout", () => {
        tooltip.style.display = "none";

        // Restore original color
        const originalFill = element.getAttribute("data-original-fill");
        if (originalFill) {
          element.setAttribute("fill", originalFill);
        }
      });
    });
  }

  // Add chart controls
  function addChartControls() {
    // Create controls container if it doesn't exist
    let controlsContainer = document.getElementById("chart-controls");
    if (!controlsContainer) {
      controlsContainer = document.createElement("div");
      controlsContainer.id = "chart-controls";
      controlsContainer.className = "chart-controls";
      chartPreview.insertAdjacentElement("afterend", controlsContainer);
    }

    // Clear existing controls
    controlsContainer.innerHTML = "";

    // Add chart parameter controls
    controlsContainer.innerHTML = `
      <div class="control-group">
        <h3>Chart Parameters</h3>
        <div class="control-row">
          <label for="chart-title">Title:</label>
          <input type="text" id="chart-title" value="${chartParameters.title}">
        </div>
        <div class="control-row">
          <label for="chart-x-label">X-Axis Label:</label>
          <input type="text" id="chart-x-label" value="${chartParameters.x_label}">
        </div>
        <div class="control-row">
          <label for="chart-y-label">Y-Axis Label:</label>
          <input type="text" id="chart-y-label" value="${chartParameters.y_label}">
        </div>
        <div class="control-row">
          <label for="chart-width">Width:</label>
          <input type="number" id="chart-width" value="${chartParameters.width}" min="300" max="1200">
        </div>
        <div class="control-row">
          <label for="chart-height">Height:</label>
          <input type="number" id="chart-height" value="${chartParameters.height}" min="200" max="800">
        </div>
        <button id="update-chart-btn">Update Chart</button>
      </div>
      <div class="control-group">
        <h3>Actions</h3>
        <button id="download-svg-btn">Download SVG</button>
        <button id="download-png-btn">Download PNG</button>
        <button id="list-data-files-btn">Browse Data Files</button>
      </div>
    `;

    // Add event listeners to controls
    document
      .getElementById("update-chart-btn")
      .addEventListener("click", () => {
        // Update chart parameters
        chartParameters.title = document.getElementById("chart-title").value;
        chartParameters.x_label =
          document.getElementById("chart-x-label").value;
        chartParameters.y_label =
          document.getElementById("chart-y-label").value;
        chartParameters.width = parseInt(
          document.getElementById("chart-width").value
        );
        chartParameters.height = parseInt(
          document.getElementById("chart-height").value
        );

        // Regenerate chart with new parameters
        const chartType = document.querySelector(
          'input[name="chart-type"]:checked'
        ).value;
        generateChart(currentDataFile, chartType);
      });

    // Download SVG button
    document
      .getElementById("download-svg-btn")
      .addEventListener("click", () => {
        if (!currentSvg) {
          showNotification("No chart available to download", "error");
          return;
        }

        // Create a blob from the SVG content
        const blob = new Blob([currentSvg], { type: "image/svg+xml" });
        const url = URL.createObjectURL(blob);

        // Create download link
        const a = document.createElement("a");
        a.href = url;
        a.download = `chart_${new Date().getTime()}.svg`;
        document.body.appendChild(a);
        a.click();

        // Clean up
        setTimeout(() => {
          document.body.removeChild(a);
          URL.revokeObjectURL(url);
        }, 100);
      });

    // Download PNG button
    document
      .getElementById("download-png-btn")
      .addEventListener("click", () => {
        if (!currentSvg) {
          showNotification("No chart available to download", "error");
          return;
        }

        // Create a canvas element
        const canvas = document.createElement("canvas");
        const ctx = canvas.getContext("2d");

        // Create an image from the SVG
        const img = new Image();
        const svgBlob = new Blob([currentSvg], { type: "image/svg+xml" });
        const url = URL.createObjectURL(svgBlob);

        img.onload = function () {
          // Set canvas dimensions
          canvas.width = img.width;
          canvas.height = img.height;

          // Draw image to canvas
          ctx.drawImage(img, 0, 0);

          // Convert to PNG
          try {
            const pngUrl = canvas.toDataURL("image/png");

            // Create download link
            const a = document.createElement("a");
            a.href = pngUrl;
            a.download = `chart_${new Date().getTime()}.png`;
            document.body.appendChild(a);
            a.click();

            // Clean up
            setTimeout(() => {
              document.body.removeChild(a);
            }, 100);
          } catch (e) {
            console.error("Error converting to PNG:", e);
            showNotification(
              "Error creating PNG. Try downloading as SVG instead.",
              "error"
            );
          }

          // Clean up
          URL.revokeObjectURL(url);
        };

        img.src = url;
      });

    // List data files button
    document
      .getElementById("list-data-files-btn")
      .addEventListener("click", () => {
        listDataFiles();
      });
  }

  // List available data files
  async function listDataFiles() {
    try {
      const response = await fetch("/api/data-files");

      if (!response.ok) {
        throw new Error("Failed to fetch data files");
      }

      const result = await response.json();

      if (result.files && result.files.length > 0) {
        showDataFilesModal(result.files);
      } else {
        showNotification("No data files available", "info");
      }
    } catch (error) {
      console.error("Error listing data files:", error);
      showNotification("Error fetching data files", "error");
    }
  }

  // Show data files modal
  function showDataFilesModal(files) {
    // Create modal container if it doesn't exist
    let modal = document.getElementById("data-files-modal");
    if (!modal) {
      modal = document.createElement("div");
      modal.id = "data-files-modal";
      modal.className = "modal";
      document.body.appendChild(modal);
    }

    // Generate file list HTML
    const fileListHtml = files
      .map(
        (file) => `
      <div class="data-file-item">
        <span class="file-name">${file.name}</span>
        <span class="file-size">${formatFileSize(file.size)}</span>
        <span class="file-date">${file.modified}</span>
        <button class="load-file-btn" data-filename="${file.name}">Load</button>
      </div>
    `
      )
      .join("");

    // Set modal content
    modal.innerHTML = `
      <div class="modal-content">
        <div class="modal-header">
          <h2>Available Data Files</h2>
          <span class="close-modal">&times;</span>
        </div>
        <div class="modal-body">
          <div class="data-files-list">
            <div class="data-file-header">
              <span>Filename</span>
              <span>Size</span>
              <span>Modified</span>
              <span>Action</span>
            </div>
            ${fileListHtml}
          </div>
        </div>
      </div>
    `;

    // Show modal
    modal.style.display = "block";

    // Add event listeners
    modal.querySelector(".close-modal").addEventListener("click", () => {
      modal.style.display = "none";
    });

    // Close modal when clicking outside
    window.addEventListener("click", (event) => {
      if (event.target === modal) {
        modal.style.display = "none";
      }
    });

    // Add load file button event listeners
    const loadButtons = modal.querySelectorAll(".load-file-btn");
    loadButtons.forEach((button) => {
      button.addEventListener("click", () => {
        const filename = button.getAttribute("data-filename");
        currentDataFile = filename;

        // Close modal
        modal.style.display = "none";

        // Generate chart with selected file
        const chartType = document.querySelector(
          'input[name="chart-type"]:checked'
        ).value;
        generateChart(filename, chartType);
      });
    });
  }

  // Show loading indicator
  function showLoadingIndicator(container) {
    container.innerHTML = `
      <div class="loading-indicator">
        <div class="spinner"></div>
        <p>Loading chart...</p>
      </div>
    `;
  }

  // Show notification
  function showNotification(message, type = "info") {
    // Create notification container if it doesn't exist
    let notificationContainer = document.getElementById(
      "notification-container"
    );
    if (!notificationContainer) {
      notificationContainer = document.createElement("div");
      notificationContainer.id = "notification-container";
      document.body.appendChild(notificationContainer);
    }

    // Create notification element
    const notification = document.createElement("div");
    notification.className = `notification ${type}`;
    notification.innerHTML = `
      <span class="notification-message">${message}</span>
      <span class="notification-close">&times;</span>
    `;

    // Add to container
    notificationContainer.appendChild(notification);

    // Add close button event listener
    notification
      .querySelector(".notification-close")
      .addEventListener("click", () => {
        notification.classList.add("fade-out");
        setTimeout(() => {
          notificationContainer.removeChild(notification);
        }, 300);
      });

    // Auto-remove after 5 seconds
    setTimeout(() => {
      if (notification.parentNode === notificationContainer) {
        notification.classList.add("fade-out");
        setTimeout(() => {
          if (notification.parentNode === notificationContainer) {
            notificationContainer.removeChild(notification);
          }
        }, 300);
      }
    }, 5000);
  }

  // Add CSS for new components
  function addStyles() {
    const styleElement = document.createElement("style");
    styleElement.textContent = `
      /* Chart Controls */
      .chart-controls {
        margin-top: 20px;
        padding: 15px;
        background-color: #f5f5f5;
        border-radius: 5px;
        box-shadow: 0 2px 4px rgba(0,0,0,0.1);
      }
      
      .control-group {
        margin-bottom: 15px;
      }
      
      .control-group h3 {
        margin-top: 0;
        margin-bottom: 10px;
        font-size: 16px;
        color: #333;
      }
      
      .control-row {
        display: flex;
        align-items: center;
        margin-bottom: 8px;
      }
      
      .control-row label {
        width: 100px;
        font-weight: 500;
      }
      
      .control-row input {
        flex: 1;
        padding: 6px 8px;
        border: 1px solid #ccc;
        border-radius: 4px;
      }
      
      button {
        padding: 8px 12px;
        background-color: #4a90e2;
        color: white;
        border: none;
        border-radius: 4px;
        cursor: pointer;
        margin-right: 8px;
        margin-top: 8px;
      }
      
      button:hover {
        background-color: #3a80d2;
      }
      
      /* Modal */
      .modal {
        display: none;
        position: fixed;
        z-index: 1000;
        left: 0;
        top: 0;
        width: 100%;
        height: 100%;
        background-color: rgba(0,0,0,0.5);
      }
      
      .modal-content {
        background-color: white;
        margin: 10% auto;
        padding: 20px;
        width: 80%;
        max-width: 700px;
        border-radius: 5px;
        box-shadow: 0 4px 8px rgba(0,0,0,0.2);
      }
      
      .modal-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 15px;
      }
      
      .close-modal {
        font-size: 24px;
        font-weight: bold;
        cursor: pointer;
      }
      
      .data-files-list {
        max-height: 400px;
        overflow-y: auto;
      }
      
      .data-file-header {
        display: grid;
        grid-template-columns: 2fr 1fr 1fr 1fr;
        padding: 10px;
        background-color: #f0f0f0;
        font-weight: bold;
        border-bottom: 1px solid #ddd;
      }
      
      .data-file-item {
        display: grid;
        grid-template-columns: 2fr 1fr 1fr 1fr;
        padding: 10px;
        border-bottom: 1px solid #eee;
      }
      
      .data-file-item:hover {
        background-color: #f9f9f9;
      }
      
      /* Notifications */
      #notification-container {
        position: fixed;
        top: 20px;
        right: 20px;
        z-index: 1000;
      }
      
      .notification {
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 12px 15px;
        margin-bottom: 10px;
        border-radius: 4px;
        box-shadow: 0 2px 5px rgba(0,0,0,0.2);
        min-width: 250px;
        max-width: 400px;
        animation: slide-in 0.3s ease-out;
      }
      
      .notification.info {
        background-color: #e3f2fd;
        border-left: 4px solid #2196f3;
      }
      
      .notification.success {
        background-color: #e8f5e9;
        border-left: 4px solid #4caf50;
      }
      
      .notification.error {
        background-color: #ffebee;
        border-left: 4px solid #f44336;
      }
      
      .notification-close {
        cursor: pointer;
        font-weight: bold;
        margin-left: 10px;
      }
      
      .notification.fade-out {
        animation: fade-out 0.3s ease-out forwards;
      }
      
      /* Loading indicator */
      .loading-indicator {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        padding: 30px;
      }
      
      .spinner {
        width: 40px;
        height: 40px;
        border: 4px solid rgba(0,0,0,0.1);
        border-radius: 50%;
        border-top-color: #4a90e2;
        animation: spin 1s linear infinite;
      }
      
      /* Animations */
      @keyframes spin {
        to { transform: rotate(360deg); }
      }
      
      @keyframes slide-in {
        from { transform: translateX(100%); opacity: 0; }
        to { transform: translateX(0); opacity: 1; }
      }
      
      @keyframes fade-out {
        from { opacity: 1; }
        to { opacity: 0; }
      }
    `;

    document.head.appendChild(styleElement);
  }

  // Add styles when DOM is loaded
  addStyles();

  // Initialize SVG chart interactivity if SVG already exists
  const existingSvg = chartPreview.querySelector("svg");
  if (existingSvg) {
    enhanceSvgInteractivity(existingSvg);
  }
});

/**
 * SVG Manipulation Utilities
 */
const SvgUtils = {
  // Add data attributes to SVG elements for better interactivity
  enhanceSvgData: function (svgElement, data) {
    if (!svgElement || !data) return;

    // For line charts - add data to circles
    const circles = svgElement.querySelectorAll("circle");
    if (circles.length > 0 && data.length >= circles.length) {
      circles.forEach((circle, index) => {
        if (data[index]) {
          circle.setAttribute("data-value", data[index].value);
          circle.setAttribute("data-category", data[index].category);
        }
      });
    }

    // For bar charts - add data to rectangles
    const rects = svgElement.querySelectorAll("rect.bar");
    if (rects.length > 0 && data.length >= rects.length) {
      rects.forEach((rect, index) => {
        if (data[index]) {
          rect.setAttribute("data-value", data[index].value);
          rect.setAttribute("data-category", data[index].category);
        }
      });
    }
  },

  // Add zoom functionality to SVG
  addZoomControls: function (svgElement) {
    if (!svgElement) return;

    // Create zoom controls container
    const zoomControls = document.createElement("div");
    zoomControls.className = "svg-zoom-controls";
    zoomControls.innerHTML = `
      <button class="zoom-in">+</button>
      <button class="zoom-out">-</button>
      <button class="zoom-reset">Reset</button>
    `;

    // Insert controls before the SVG
    svgElement.parentNode.insertBefore(zoomControls, svgElement);

    // Current zoom level
    let zoomLevel = 1;
    const zoomStep = 0.1;
    const maxZoom = 3;
    const minZoom = 0.5;

    // Get the SVG viewBox
    const originalViewBox = svgElement.getAttribute("viewBox");

    // Add event listeners
    zoomControls.querySelector(".zoom-in").addEventListener("click", () => {
      if (zoomLevel < maxZoom) {
        zoomLevel += zoomStep;
        applyZoom();
      }
    });

    zoomControls.querySelector(".zoom-out").addEventListener("click", () => {
      if (zoomLevel > minZoom) {
        zoomLevel -= zoomStep;
        applyZoom();
      }
    });

    zoomControls.querySelector(".zoom-reset").addEventListener("click", () => {
      zoomLevel = 1;
      svgElement.setAttribute("viewBox", originalViewBox);
      svgElement.style.transform = "none";
    });

    // Apply zoom function
    function applyZoom() {
      svgElement.style.transform = `scale(${zoomLevel})`;
      svgElement.style.transformOrigin = "center center";
    }
  },

  // Export SVG to various formats
  exportSvg: {
    toSvgString: function (svgElement) {
      if (!svgElement) return null;
      return new XMLSerializer().serializeToString(svgElement);
    },

    toDataUrl: function (svgElement) {
      if (!svgElement) return null;
      const svgString = this.toSvgString(svgElement);
      return (
        "data:image/svg+xml;charset=utf-8," + encodeURIComponent(svgString)
      );
    },

    downloadAsSvg: function (svgElement, filename = "chart.svg") {
      if (!svgElement) return;

      const svgUrl = this.toDataUrl(svgElement);
      const downloadLink = document.createElement("a");
      downloadLink.href = svgUrl;
      downloadLink.download = filename;
      document.body.appendChild(downloadLink);
      downloadLink.click();
      document.body.removeChild(downloadLink);
    },

    toPngDataUrl: function (svgElement, callback) {
      if (!svgElement) return null;

      const svgString = this.toSvgString(svgElement);
      const svgBlob = new Blob([svgString], {
        type: "image/svg+xml;charset=utf-8",
      });
      const url = URL.createObjectURL(svgBlob);

      const img = new Image();
      img.onload = function () {
        const canvas = document.createElement("canvas");
        canvas.width = svgElement.viewBox.baseVal.width || img.width;
        canvas.height = svgElement.viewBox.baseVal.height || img.height;

        const ctx = canvas.getContext("2d");
        ctx.drawImage(img, 0, 0);

        URL.revokeObjectURL(url);

        const pngDataUrl = canvas.toDataURL("image/png");
        if (callback) callback(pngDataUrl);
      };

      img.src = url;
    },
  },
};
