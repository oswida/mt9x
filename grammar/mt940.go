package grammar

import (
	"fmt"
	"mt9x/bundle"
	"strings"
)

// Grammar for MT940 file, according standard available here:
// https://www2.swift.com/knowledgecentre/publications/us9m_20240719/2.0

// Message represents MT940 standard message structure.
type MT940Message struct {
	// Specifies the reference assigned by the Sender to unambiguously identify the message.
	TransactionRefNo string `parser:"T20 @StringX CRLF" json:"tag20"`
	// Contains the field 20 Transaction Reference Number of the request message (response to MT920 Request Message).
	RelatedReference *string `parser:"(T21 @StringX CRLF)?" json:"tag21,omitempty"`
	// Identifies the account and optionally the identifier code of the account owner for which the statement is sent.
	// Need some examples, optional
	AccountIdentification AccountIdent `parser:"(T25|T25P) @@ CRLF" json:"tag25"`
	// Contains the sequential number of the statement, optionally followed by the sequence number of the message
	// within that statement when more than one message is sent for one statement.
	StatementNumber StatementNumber `parser:"T28C @@ CRLF" json:"tag28"`
	// Specifies, for the (intermediate - M) opening balance, whether it is a debit or credit balance,
	// the date, the currency and the amount of the balance.
	OpeningBalance Balance `parser:"(T60F|T60M) @@ (CRLF|EOF)" json:"tag60"`
	// Statement information
	Statements []StatementSection `parser:"@@*" json:"statements,omitempty"`
	// Specifies, for the (intermediate) closing balance.
	ClosingBalance Balance `parser:"(T62F|T62M) @@ (CRLF|EOF)" json:"tag62"`
	// Indicates the funds which are available to the account owner (if credit balance)
	// or the balance which is subject to interest charges (if debit balance).
	ClosingAvailableBalance *Balance `parser:"(T64 @@ (CRLF|EOF))?" json:"tag64,omitempty"`
	// Indicates the funds which are available to the account owner
	// (if a credit or debit balance) for the specified forward value date.
	ForwardAvailableBalance []Balance `parser:"(T65 @@ (CRLF|EOF))*" json:"tag65,omitempty"`
	// Summarizing owner info
	AccountOwnerInfo []string `parser:"(T86 @StringX (CRLF @StringX)* (CRLF|EOF))?" json:"tag86,omitempty"`
}

// Validate validates MT940 messages according "Network Validated Rules"
func (m MT940Message) Validate() error {
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
