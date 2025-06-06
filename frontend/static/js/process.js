import { showNotification } from './notification.js';

export const process = async (form) => {
  try {
    console.log('Sending request to /upload...');

    const response = await fetch('/upload', {
      method: 'POST',
      body: form,
      headers: {
        Accept: 'application/json',
      },
    });

    if (!response.ok) {
      const errorText = await response.text();
      console.error('Server response error:', errorText);
      showNotification(
        `Error processing the document. Status: ${response.status}. Please try again.`
      );
      return null;
    }

    const contentType = response.headers.get('content-type');
    if (!contentType || !contentType.includes('application/json')) {
      const responseText = await response.text();
      console.error('Non-JSON response:', responseText);
      showNotification('Server returned invalid response format.');
      return null;
    }

    let data = await response.json();
    console.log('Parsed response data:', data);

    if (data.status === 'error') {
      console.error('Server error:', data.message);
      showNotification(data.message || 'Server error occurred');
      return null;
    }

    // Check if we have the expected structure
    if (!data.units) {
      console.error('Missing units field in response:', data);
      showNotification('Server response missing required data');
      return null;
    }

    // Return both units and data file name for chart generation
    return {
      units: data.units,
      dataFile: data.data_file || `${Date.now()}.json`, // Fallback if data_file is missing
    };
  } catch (error) {
    console.error('Network or parsing error:', error);
    showNotification('Network error: ' + error.message);
    return null;
  }
};
