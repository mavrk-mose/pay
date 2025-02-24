package util

import (
	"encoding/csv"
	"fmt"
	"os"
)

// WriteCSV is a generic function that writes any slice of items to a CSV file.
// - fileName: the output CSV file path.
// - items: slice of items of type T.
// - headers: a slice of column header names.
// - rowFunc: a function that converts an item of type T to a []string representing a row.
func WriteCSV[T any](fileName string, items []T, headers []string, rowFunc func(T) []string) error {
	// Create or overwrite the file.
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	// Write header.
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write each row by calling rowFunc for each item.
	for _, item := range items {
		row := rowFunc(item)
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	// Flush to ensure all data is written.
	writer.Flush()
	if err := writer.Error(); err != nil {
		return fmt.Errorf("error flushing CSV writer: %w", err)
	}

	return nil
}

// ReadCSV is a generic function that reads a CSV file and returns a slice of items of type T.
// - fileName: the path to the CSV file.
// - rowParser: a function that converts a CSV row ([]string) into an instance of type T.
// - skipHeader: if true, skips the first row in the CSV (assumed to be headers).
func ReadCSV[T any](fileName string, rowParser func([]string) (T, error), skipHeader bool) ([]T, error) {
	// Open the CSV file.
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read all rows from the file.
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV file: %w", err)
	}

	var results []T
	for i, row := range rows {
		// Skip header row if required.
		if skipHeader && i == 0 {
			continue
		}
		item, err := rowParser(row)
		if err != nil {
			return nil, fmt.Errorf("failed to parse row %d: %w", i, err)
		}
		results = append(results, item)
	}
	return results, nil
}

