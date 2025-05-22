document.addEventListener('DOMContentLoaded', function() {
    const chartOutput = document.getElementById('chart-output');
    const generateBtn = document.getElementById('generate-btn');
    const addRowBtn = document.getElementById('add-row-btn');
    const dataTable = document.getElementById('data-table');

    // Generate chart when button is clicked
    generateBtn.addEventListener('click', generateChart);

    // Add new row to data table
    addRowBtn.addEventListener('click', function() {
        const newRow = document.createElement('tr');
        newRow.innerHTML = `
            <td><input type="text" class="data-label" value="New"></td>
            <td><input type="number" class="data-value" value="0"></td>
            <td><button class="remove-btn">×</button></td>
        `;
        dataTable.querySelector('tbody').appendChild(newRow);
        addRemoveListener(newRow.querySelector('.remove-btn'));
    });

    // Add event listeners to existing remove buttons
    document.querySelectorAll('.remove-btn').forEach(btn => {
        addRemoveListener(btn);
    });

    // Generate initial chart
    generateChart();

    function addRemoveListener(btn) {
        btn.addEventListener('click', function() {
            if (dataTable.querySelectorAll('tbody tr').length > 1) {
                btn.closest('tr').remove();
            } else {
                alert('You must have at least one data point.');
            }
        });
    }

    function generateChart() {
        const chartType = document.getElementById('chart-type').value;
        const title = document.getElementById('chart-title').value;
        const xLabel = document.getElementById('x-label').value;
        const yLabel = document.getElementById('y-label').value;
        const width = parseInt(document.getElementById('width').value);
        const height = parseInt(document.getElementById('height').value);
        const showGrid = document.getElementById('show-grid').checked;

        // Collect data from table
        const data = [];
        document.querySelectorAll('#data-table tbody tr').forEach(row => {
            const label = row.querySelector('.data-label').value;
            const value = parseFloat(row.querySelector('.data-value').value);
            if (label && !isNaN(value)) {
                data.push({ label, value });
            }
        });

        // Prepare request data
        const requestData = {
            chartType,
            data,
            options: {
                title,
                xLabel,
                yLabel,
                width,
                height,
                showGrid
            }
        };

        // Show loading state
        chartOutput.innerHTML = '<p>Generating chart...</p>';

        // Send request to server
        fetch('/api/generate-chart', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(requestData)
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Chart generation failed');
            }
            return response.text();
        })
        .then(svg => {
            chartOutput.innerHTML = svg;
        })
        .catch(error => {
            chartOutput.innerHTML = `<p class="error">Error: ${error.message}</p>`;
            console.error('Error:', error);
        });
    }
});
