package util

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

/*
ParseDocumentToJSON reads a text document, extracts numeric data with units, and saves the extracted information as a formatted JSON file.

It reads the input file in 512-byte chunks, processes each chunk with GetData(), and collects all extracted DataPoint entries.

Parameters:
  filePath string   - path to the input text file
  outputPath string - path where the output JSON file should be saved

Returns:
  error - if any file operation or JSON encoding fails, the error is returned.

Dependencies:
  - GetData(text string) []DataPoint: used to extract data from text chunks

Example:
  err := ParseDocumentToJSON("report.txt", "data.json")
  data.json will contain:
  [
    { "value": 23.5, "unit": "°C" },
    { "value": 120.0, "unit": "vehicles/hr" }
  ]
*/

func ParseDocumentToJSON(filePath string, outputPath string) error {
	const chunkSize = 512
	var allData []DataPoint

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	buffer := make([]byte, chunkSize)

	for {
		n, err := reader.Read(buffer)
		if n > 0 {
			chunkText := string(buffer[:n])
			data := GetData(chunkText)
			allData = append(allData, data...)
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading file: %w", err)
		}
	}

	// Write to JSON file
	jsonFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer jsonFile.Close()

	encoder := json.NewEncoder(jsonFile)
	encoder.SetIndent("", "  ")

	err = encoder.Encode(allData)
	if err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}
