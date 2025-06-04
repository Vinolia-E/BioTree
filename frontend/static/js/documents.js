// Function to fetch user documents
export async function fetchUserDocuments() {
  try {
    const response = await fetch('/api/data-files');
    
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    
    const data = await response.json();
    
    if (data.status === 'error') {
      throw new Error(data.message || 'Failed to fetch documents');
    }
    
    return data.files || [];
  } catch (error) {
    console.error('Error fetching documents:', error);
    return [];
  }
}

// Function to display user documents
export function displayUserDocuments(documents) {
  const documentsList = document.getElementById('documents-list');
  
  if (!documents || documents.length === 0) {
    documentsList.innerHTML = '<p class="no-documents">No documents found. Upload a document to get started.</p>';
    return;
  }
  
  // Sort documents by modified date (newest first)
  documents.sort((a, b) => new Date(b.modified) - new Date(a.modified));
  
  let html = '';
  
  documents.forEach(doc => {
    const fileName = doc.name;
    const fileSize = formatFileSize(doc.size);
    const modified = new Date(doc.modified).toLocaleString();
    const units = doc.units && doc.units.length > 0 
      ? `<span class="units-count">${doc.units.length} units</span>` 
      : '<span class="units-count">No units</span>';
    
    html += `
      <div class="document-item">
        <div class="document-info">
          <h4 class="document-name">${fileName}</h4>
          <p class="document-meta">
            ${fileSize} • ${modified}<br>
            ${units}
          </p>
        </div>
        <div class="document-actions">
          <button class="view-btn" data-filename="${fileName}">View</button>
          <button class="generate-chart-btn" data-filename="${fileName}" ${doc.units.length === 0 ? 'disabled' : ''}>
            Chart
          </button>
        </div>
      </div>
    `;
  });
  
  documentsList.innerHTML = html;
  
  // Add event listeners to buttons
  addDocumentButtonListeners();
}

// Function to add event listeners to document buttons
function addDocumentButtonListeners() {
  // View content buttons
  document.querySelectorAll('.view-btn').forEach(button => {
    button.addEventListener('click', async (e) => {
      const fileName = e.target.dataset.filename;
      await viewDocumentContent(fileName);
    });
  });
  
  // Generate chart buttons
  document.querySelectorAll('.generate-chart-btn').forEach(button => {
    button.addEventListener('click', (e) => {
      const fileName = e.target.dataset.filename;
      generateChartForDocument(fileName);
    });
  });
}

// Function to view document content
async function viewDocumentContent(fileName) {
  try {
    const response = await fetch(`/data/${fileName}`);
    
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    
    const data = await response.json();
    
    // Display the content in the modal
    const modal = document.getElementById('document-content-modal');
    const modalTitle = document.getElementById('modal-title');
    const documentContent = document.getElementById('document-content');
    
    modalTitle.textContent = `Content of ${fileName}`;
    
    // Format the JSON data
    documentContent.innerHTML = `<pre>${JSON.stringify(data, null, 2)}</pre>`;
    
    // Show the modal
    modal.style.display = 'block';
    
    // Add close button functionality
    document.getElementById('close-modal').onclick = function() {
      modal.style.display = 'none';
    };
    
    // Close modal when clicking outside
    window.onclick = function(event) {
      if (event.target === modal) {
        modal.style.display = 'none';
      }
    };
    
  } catch (error) {
    console.error('Error viewing document content:', error);
    alert('Failed to load document content. Please try again.');
  }
}

// Function to generate chart for document
function generateChartForDocument(fileName) {
  // Get the controls container
  const controlsContainer = document.getElementById('controls-container');
  
  // Clear previous content
  controlsContainer.innerHTML = '';
  
  // Show the controls container
  controlsContainer.style.display = 'block';
  
  // Scroll to the controls container
  controlsContainer.scrollIntoView({ behavior: 'smooth' });
  
  // Create chart controls
  controlsContainer.insertAdjacentHTML('beforeend', `
    <div class="control-group">
      <h3>Chart for ${fileName}</h3>
      <div class="control-row">
        <label for="unit-select">Select Unit:</label>
        <select id="unit-select">
          <option value="">Loading units...</option>
        </select>
      </div>
      <div class="control-row">
        <label for="chart-title">Title:</label>
        <input type="text" id="chart-title" placeholder="Enter chart title">
      </div>
      <div class="control-row">
        <label for="chart-x-label">X-Axis Label:</label>
        <input type="text" id="chart-x-label" placeholder="Enter X-axis label">
      </div>
      <div class="control-row">
        <label for="chart-y-label">Y-Axis Label:</label>
        <input type="text" id="chart-y-label" placeholder="Enter Y-axis label">
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
      <button id="update-chart-btn" class="primary-button">Generate Chart</button>
    </div>
    <div class="control-group">
      <h3>Actions</h3>
      <div class="button-row">
        <button id="download-svg-btn" class="secondary-button" disabled>Download SVG</button>
        <button id="reset-chart-btn" class="secondary-button">Reset Chart</button>
      </div>
    </div>
    <div id="chart-container" class="chart-display-area"></div>
  `);
  
  // Load units for the document
  loadUnitsForDocument(fileName);
  
  // Add event listeners for chart generation
  document.getElementById('update-chart-btn').addEventListener('click', () => {
    updateChart(fileName);
  });
  
  // Add event listener for download button
  document.getElementById('download-svg-btn').addEventListener('click', downloadSVG);
  
  // Add event listener for reset button
  document.getElementById('reset-chart-btn').addEventListener('click', () => {
    document.getElementById('chart-container').innerHTML = '';
    document.getElementById('download-svg-btn').disabled = true;
  });
}

// Function to load units for a document
async function loadUnitsForDocument(fileName) {
  try {
    const response = await fetch('/api/data-files');
    
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    
    const data = await response.json();
    
    if (data.status === 'error') {
      throw new Error(data.message || 'Failed to fetch documents');
    }
    
    const fileData = data.files.find(file => file.name === fileName);
    
    if (!fileData) {
      throw new Error('Document not found');
    }
    
    const unitSelect = document.getElementById('unit-select');
    unitSelect.innerHTML = '';
    
    // Add "All" option
    const allOption = document.createElement('option');
    allOption.value = '';
    allOption.textContent = 'All Units';
    unitSelect.appendChild(allOption);
    
    // Add each unit as an option
    fileData.units.forEach(unit => {
      const option = document.createElement('option');
      option.value = unit;
      option.textContent = unit;
      unitSelect.appendChild(option);
    });
    
  } catch (error) {
    console.error('Error loading units:', error);
    document.getElementById('unit-select').innerHTML = '<option value="">Failed to load units</option>';
  }
}

// Function to update chart
async function updateChart(fileName) {
  const updateChartBtn = document.getElementById('update-chart-btn');
  const originalText = updateChartBtn.textContent;
  updateChartBtn.textContent = 'Generating...';
  updateChartBtn.disabled = true;
  
  try {
    const chartRequest = {
      dataFile: fileName,
      chartType: document.getElementById('chart-type').value,
      unit: document.getElementById('unit-select').value,
      title: document.getElementById('chart-title').value,
      xLabel: document.getElementById('chart-x-label').value,
      yLabel: document.getElementById('chart-y-label').value,
      width: parseInt(document.getElementById('chart-width').value),
      height: parseInt(document.getElementById('chart-height').value)
    };
    
    console.log('Sending chart request:', chartRequest);
    
    // Use the API endpoint instead of the legacy endpoint
    const response = await fetch('/api/generate-chart', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json'
      },
      body: JSON.stringify(chartRequest)
    });
    
    if (!response.ok) {
      const errorText = await response.text();
      console.error('Server error response:', errorText);
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    
    const result = await response.json();
    
    if (result.status === 'error') {
      throw new Error(result.message || 'Failed to generate chart');
    }
    
    // Display the chart
    const chartContainer = document.getElementById('chart-container');
    chartContainer.innerHTML = result.svg;
    
    // Enable download button
    document.getElementById('download-svg-btn').disabled = false;
    
    // Store SVG for download
    window.currentSVG = result.svg;
    
  } catch (error) {
    console.error('Error generating chart:', error);
    alert('Error generating chart: ' + error.message);
  } finally {
    // Reset button state
    updateChartBtn.textContent = originalText;
    updateChartBtn.disabled = false;
  }
}

// Function to download SVG
function downloadSVG() {
  if (!window.currentSVG) {
    alert('No chart available to download.');
    return;
  }
  
  const blob = new Blob([window.currentSVG], { type: 'image/svg+xml' });
  const url = URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = 'chart.svg';
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  URL.revokeObjectURL(url);
}

// Helper function to format file size
function formatFileSize(bytes) {
  if (bytes < 1024) {
    return bytes + ' B';
  } else if (bytes < 1024 * 1024) {
    return (bytes / 1024).toFixed(1) + ' KB';
  } else {
    return (bytes / (1024 * 1024)).toFixed(1) + ' MB';
  }
}
