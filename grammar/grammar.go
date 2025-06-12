package grammar

import (
	"fmt"
	"mt9x/bundle"
	"mt9x/parser"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"
	// "github.com/alecthomas/participle/v2/lexer"
)

type StatementSection struct {
	// Contains the details of each transaction.
	Statement Statement `parser:"T61 @@ (CRLF|EOF)" json:"tag61"`
	// Contains additional information about the transaction detailed in the preceding statement line
	// and which is to be passed on to the account owner.
	AccountOwnerInfo []string `parser:"(T86 @StringX (CRLF @StringX)* (CRLF|EOF))?" json:"tag86,omitempty"`
}

type AccountIdent struct {
	Account   string  `parser:"@StringX" json:"account"`
	IdentCode *string `parser:"(CRLF @StringX)?" json:"ident_code,omitempty"` //TODO: validate ident code 4!a2!a2!c[3!c]
}

type StatementNumber struct {
	StatementNo string  `parser:"@Number" json:"stmt_number"`
	SequenceNo  *string `parser:"(Slash @Number)?" json:"seq_number,omitempty"`
}

type Balance struct {
	DCMark   string              `parser:"@DCMark" json:"dc_mark"`
	Date     parser.SixDigitDate `parser:"@Date" json:"date"`
	Currency string              `parser:"@Currency" json:"currency"`
	Amount   parser.CommaDecimal `parser:"@Amount" json:"amount"`
}

type Statement struct {
	ValueDate            parser.SixDigitDate   `parser:"@Date" json:"value_date"`
	EntryDate            *parser.FourDigitDate `parser:"@EntryDate?" json:"entry_date,omitempty"`
	DCMark               string                `parser:"@RDCMark" json:"dc_mark"`
	FundsCode            *string               `parser:"@BigLetter?" json:"funds_code,omitempty"`
	Amount               parser.CommaDecimal   `parser:"@Amount" json:"amount"`
	TransactionIdent     string                `parser:"@TransIdent" json:"trx_ident"`
	Reference            string                `parser:"@StringXNo2Slash" json:"owner_ref"`
	InstitutionReference *string               `parser:"(TwoSlashes @StringXNo2Slash?)?" json:"institution_ref,omitempty"`
	Details              *string               `parser:"(CRLF @StringX)?" json:"details,omitempty"`
}

// --- VALIDATIONS ---

// Validate validates balance field according "Network Validated Rules"
func (b *Balance) Validate(cp *bundle.CurrencyProvider) error {
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
	if !isCorrectTransactionIdent(s.TransactionIdent, sicp) {
		return fmt.Errorf("bad transaction ident: %s", s.TransactionIdent)
	}

	return nil
}
