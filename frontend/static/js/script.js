document.addEventListener('DOMContentLoaded', () => {
  // DOM Elements
  const dropArea = document.getElementById('dropArea');
  const fileInput = document.getElementById('document-upload');
  const filePreview = document.getElementById('filePreview');
  const chartPreview = document.getElementById('chartPreview');
  const processBtn = document.getElementById('processBtn');
  const chartOptions = document.querySelectorAll('input[name="chart-type"]');

  // Accepted files
  const fileTypes = [
    'application/pdf',
    'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
    'text/plain',
  ];

  // Event Listeners for drag and drop
  ['dragenter', 'dragover', 'dragleave', 'drop'].forEach((eventName) => {
    dropArea.addEventListener(eventName, preventDefaults, false);
  });

  function preventDefaults(e) {
    e.preventDefault();
    e.stopPropagation();
  }

  ['dragenter', 'dragover'].forEach((eventName) => {
    dropArea.addEventListener(eventName, highlight, false);
  });

  ['dragleave', 'drop'].forEach((eventName) => {
    dropArea.addEventListener(eventName, unhighlight, false);
  });

  function highlight() {
    dropArea.classList.add('highlight');
  }

  function unhighlight() {
    dropArea.classList.remove('highlight');
  }

  // Handle dropped files
  dropArea.addEventListener('drop', handleDrop, false);

  function handleDrop(e) {
    const dt = e.dataTransfer;
    const files = dt.files;

    if (files.length) {
      fileInput.files = files;
      handleFiles(files);
    }
  }

  // Handle file selection via input
  fileInput.addEventListener('change', function () {
    handleFiles(this.files);
  });

  // Process the selected files
  function handleFiles(files) {
    const file = files[0];

    if (validFileType(file)) {
      displayFileInfo(file);
      processBtn.disabled = false;
    } else {
      filePreview.innerHTML = `
        <p class="error">Error: ${file.name} is not a supported file type. 
        Please select a PDF, DOCX, or TXT file.</p>
      `;
      processBtn.disabled = true;
    }
  }

  // Check if file type is valid
  function validFileType(file) {
    return (
      fileTypes.includes(file.type) ||
      file.name.endsWith('.pdf') ||
      file.name.endsWith('.docx') ||
      file.name.endsWith('.txt')
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

    // Update chart preview with placeholder
    updateChartPreview();
  }

  // Format file size
  function formatFileSize(bytes) {
    if (bytes < 1024) {
      return bytes + ' bytes';
    } else if (bytes < 1048576) {
      return (bytes / 1024).toFixed(1) + ' KB';
    } else {
      return (bytes / 1048576).toFixed(1) + ' MB';
    }
  }

  // Get file icon type
  function getFileIconType(file) {
    if (file.type === 'application/pdf' || file.name.endsWith('.pdf')) {
      return 'pdf-icon';
    } else if (
      file.type ===
        'application/vnd.openxmlformats-officedocument.wordprocessingml.document' ||
      file.name.endsWith('.docx')
    ) {
      return 'docx-icon';
    } else {
      return 'txt-icon';
    }
  }

  // Update chart preview based on selected chart type
  function updateChartPreview() {
    const selectedChart = document.querySelector(
      'input[name="chart-type"]:checked'
    ).value;

    // This would generate an actual chart but we'll just show a placeholder message
    chartPreview.innerHTML = `
      <div class="chart-placeholder ${selectedChart}-placeholder">
        <p>Preview of ${selectedChart} chart</p>
        <p>This is where the ${selectedChart} chart will be rendered</p>
      </div>
    `;
  }

  // Listen for chart type changes
  chartOptions.forEach((option) => {
    option.addEventListener('change', updateChartPreview);
  });

  // Process button click handler
  processBtn.addEventListener('click', () => {
    const file = fileInput.files[0];
    const chartType = document.querySelector(
      'input[name="chart-type"]:checked'
    ).value;

    // This should send the file to the backend but we'll just show an alert
    alert(`Processing ${file.name} with ${chartType} chart visualization.
    In a real application, this would send the file to the backend for processing.`);
  });
});
