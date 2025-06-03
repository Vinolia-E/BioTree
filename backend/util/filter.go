package util

import (
	"encoding/json"
	"fmt"
	"os"
)

// FilterAndRewriteByUnit reads a JSON file, removes entries with "(none)" unit,
// and writes the filtered data back in original flat format.
func FilterAndRewriteByUnit(filePath string) error {
	// Step 1: Read JSON from file
	dataBytes, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	var data []DataPoint
	if err := json.Unmarshal(dataBytes, &data); err != nil {
		return fmt.Errorf("unmarshaling JSON: %w", err)
	}

	// Step 2: Group by unit, skip "(none)"
	grouped := make(map[string][]DataPoint)
	for _, dp := range data {
		if dp.Unit == "(none)" {
			continue
		}
		grouped[dp.Unit] = append(grouped[dp.Unit], dp)
	}

	// Step 3: Flatten grouped data
	var filtered []DataPoint
	for _, group := range grouped {
		filtered = append(filtered, group...)
	}

	// Step 4: Write back to the same file
	outputBytes, err := json.MarshalIndent(filtered, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling JSON: %w", err)
	}

	if err := os.WriteFile(filePath, outputBytes, 0o644); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	return nil
}
