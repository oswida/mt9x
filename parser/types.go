package parser

import (
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

// MT9xMessage
type MT9xMessage interface {
	Validate() error
}

// CommaDecimal captures decimal with comma (instead of dot) as a decimal sign
type CommaDecimal struct {
	decimal.Decimal
}

func (d *CommaDecimal) Capture(values []string) error {
	if len(values) != 1 {
		return fmt.Errorf("bad capture length for CommaDecimal: %v", values)
	}
	v, err := decimal.NewFromString(strings.ReplaceAll(values[0], ",", "."))
	if err != nil {
		return err
	}

	d.Decimal = v
	return nil
}

// SixDigitDate captures dates in YYMMDD format.
type SixDigitDate struct {
	time.Time
}

func (d *SixDigitDate) Capture(values []string) error {
	if len(values) != 1 {
		return fmt.Errorf("bad capture length for SixDigitDate: %v", values)
	}
	v, err := time.Parse("060102", values[0])
	if err != nil {
		return err
	}

	d.Time = v
	return nil
}

// FourDigitDate captures dates in MMDD format.
type FourDigitDate struct {
	time.Time
}

func (d *FourDigitDate) Capture(values []string) error {
	if len(values) != 1 {
		return fmt.Errorf("bad capture length for FourDigitDate: %v", values)
	}
	v, err := time.Parse("0102", values[0])
	if err != nil {
		return err
	}

	d.Time = v
	return nil
}
