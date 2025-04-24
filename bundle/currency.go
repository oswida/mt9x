package bundle

import (
	"encoding/xml"
	"fmt"
	"strconv"
)

type entry struct {
	Currency    string `xml:"Ccy"`
	CurrencyNbr int    `xml:"CcyNbr"`
	MinorUnits  units  `xml:"CcyMnrUnts"`
}

type units uint

func (m *units) UnmarshalText(b []byte) error {
	newInt, err := strconv.ParseUint(string(b), 10, 0)
	if err == nil {
		*m = units(newInt)
	} else {
		*m = 0
	}

	return nil
}

type CurrencyProvider struct {
	Entries []entry `xml:"CcyTbl>CcyNtry,name"`
}

// NewCurrencyProvider creates new currency data provider.
func NewCurrencyProvider() (*CurrencyProvider, error) {
	cp := &CurrencyProvider{
		Entries: []entry{},
	}
	if err := cp.Load(); err != nil {
		return nil, fmt.Errorf("failed to load data: %w", err)
	}

	return cp, nil
}

// Load loads data from embedded file.
func (cp *CurrencyProvider) Load() error {
	file, err := EmbedFS.Open("resources/iso4217.xml")
	if err != nil {
		return fmt.Errorf("failed to open currency data: %w", err)
	}
	decoder := xml.NewDecoder(file)
	err = decoder.Decode(cp)
	if err != nil {
		return fmt.Errorf("failed to parse currency data: %w", err)
	}

	return nil
}

// List provides a list of available currency codes for ISO4217.
func (cp *CurrencyProvider) List() []string {
	result := make([]string, len(cp.Entries))
	for i, e := range cp.Entries {
		result[i] = e.Currency
	}

	return result
}
