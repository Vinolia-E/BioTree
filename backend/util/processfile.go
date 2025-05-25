package util

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

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
