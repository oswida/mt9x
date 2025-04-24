package bundle

import (
	"encoding/csv"
	"fmt"
)

type StatementIdentCodeProvider struct {
	Codes map[string]string
}

// NewStatementIdentificationCodeProvider creates new MT940 statement identification code provider.
func NewStatementIdentificationCodeProvider() (*StatementIdentCodeProvider, error) {
	result := &StatementIdentCodeProvider{
		Codes: make(map[string]string),
	}
	if err := result.Load(); err != nil {
		return nil, fmt.Errorf("error loading statement identification codes: %w", err)
	}

	return result, nil
}

// Load loads data from embedded file.
func (cp *StatementIdentCodeProvider) Load() error {
	file, err := EmbedFS.Open("resources/mt940sic.csv")
	if err != nil {
		return fmt.Errorf("failed to open identification data: %w", err)
	}
	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to parse identification data: %w", err)
	}
	for _, record := range records {
		cp.Codes[record[0]] = record[1]
	}

	return nil
}

// IsProperCode checks if provided statement identification code is proper.
func (cp *StatementIdentCodeProvider) IsProperCode(code string) bool {
	_, ok := cp.Codes[code]
	return ok
}
