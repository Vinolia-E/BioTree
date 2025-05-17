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
});
