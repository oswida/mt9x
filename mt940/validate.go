package mt940

import (
	"fmt"
	"mt9x/bundle"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// Validate validates MT940 messages according "Network Validated Rules"
func (m *Message) Validate() error {
	cp, err := bundle.NewCurrencyProvider()
	if err != nil {
		return fmt.Errorf("cannot create currency provider: %v", err)
	}
	sicp, err := bundle.NewStatementIdentificationCodeProvider()
	if err != nil {
		return fmt.Errorf("cannot create statement identification provider: %v", err)
	}

	if !isCorrectReference(m.TransactionRefNo) {
		return fmt.Errorf("bad transaction reference number: %s", m.TransactionRefNo)
	}

	if m.RelatedReference != nil {
		if !isCorrectReference(*m.RelatedReference) {
			return fmt.Errorf("bad related reference number: %s", *m.RelatedReference)
		}
	}

	for _, line := range m.Statements {
		err := line.Validate(sicp)
		if err != nil {
			return fmt.Errorf("error parsing MT940 statement: %w", err)
		}
	}

	if err := m.OpeningBalance.Validate(cp); err != nil {
		return fmt.Errorf("bad opening balance: %w", err)
	}

	if err := m.ClosingBalance.Validate(cp); err != nil {
		return fmt.Errorf("bad closing balance: %w", err)
	}

	if m.ClosingAvailableBalance != nil {
		if err := m.ClosingAvailableBalance.Validate(cp); err != nil {
			return fmt.Errorf("bad closing available balance: %w", err)
		}
	}

	return nil
}

// isCorrectReference checks if reference number is proper according the standard.
func isCorrectReference(ref string) bool {
	return !strings.HasPrefix(ref, "/") &&
		!strings.HasSuffix(ref, "/") &&
		!strings.Contains(ref, "//")
}

// Validate validates balance field according "Network Validated Rules"
func (b *Balance) Validate(cp *bundle.CurrencyProvider) error {
	if _, err := time.Parse("060102", b.Date); err != nil {
		return fmt.Errorf("bad date: %w", err)
	}
	if !slices.Contains(cp.List(), b.Currency) {
		return fmt.Errorf("bad currency code: %s", b.Currency)
	}
	// Amount is verified by a lexer

	return nil
}

// trimFirstRune removes first rune from the string and returns the result.
func trimFirstRune(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}

// isCorrectTransactionIdent checks if transaction identification data is proper according the standard.
func isCorrectTransactionIdent(ti string, sicp *bundle.StatementIdentCodeProvider) bool {
	if !strings.HasPrefix(ti, "S") && !strings.HasPrefix(ti, "N") && !strings.HasPrefix(ti, "F") {
		return false
	}
	if strings.HasPrefix(ti, "N") || strings.HasPrefix(ti, "F") {
		return sicp.IsProperCode(trimFirstRune(ti))
	}
	if strings.HasPrefix(ti, "S") {
		_, err := strconv.Atoi(trimFirstRune(ti))
		if err == nil {
			return true
		}
	}
	return false
}

// Validate validates statement section.
func (ss *StatementSection) Validate(sicp *bundle.StatementIdentCodeProvider) error {
	return ss.Statement.Validate(sicp)
}

// Validate validates single statement line according "Network Validated Rules".
func (s *Statement) Validate(sicp *bundle.StatementIdentCodeProvider) error {
	if _, err := time.Parse("060102", s.ValueDate); err != nil {
		return fmt.Errorf("bad value date: %w", err)
	}
	if s.EntryDate != nil {
		if _, err := time.Parse("0102", *s.EntryDate); err != nil {
			return fmt.Errorf("bad entry date: %w", err)
		}
	}
	if !isCorrectTransactionIdent(s.TransactionIdent, sicp) {
		return fmt.Errorf("bad transaction ident: %s", s.TransactionIdent)
	}

	return nil
}
