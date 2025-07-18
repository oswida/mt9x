package grammar

import (
	"fmt"
	"strings"
	"time"

	"github.com/oswida/mt9x/bundle"
)

const (
	CSVHeader = `TransactionRefNo,RelatedReference,Account,IdentCode,StmtNo,SeqNo,OB_DC,OB_Date,OB_Curr,OB_Amount,ValueDate,EntryDate,DC,FCode,Amount,TrxIdent,Reference,InstitutionRef,Details,AccOwnerInfo,CB_DC,CB_Date,CB_Curr,CB_Amount,CAB_DC,CAB_Date,CAB_Curr,CAB_Amount,FAB_DC,FAB_Date,FAB_Curr,FAB_Amount,MsgAccOwnerInfo`
)

// Grammar for MT940 file, according standard available here:
// https://www2.swift.com/knowledgecentre/publications/us9m_20240719/2.0

// Message represents MT940 standard message structure.
type MT940Message struct {
	// Specifies the reference assigned by the Sender to unambiguously identify the message.
	TransactionRefNo string `parser:"T20 @CharXSeqSlashRestrict CRLF" json:"tag20"`
	// If the MT 940 is sent in response to an MT 920 Request Message, this field must contain the field 20 Transaction Reference Number of the request message.
	RelatedReference *string `parser:"(T21 @CharXSeqSlashRestrict CRLF)?" json:"tag21,omitempty"`
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
	AccountOwnerInfo []string `parser:"(T86 @CharXSeq ((CRLF @CharXSeq?)*|EOF))?" json:"tag86,omitempty"`
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

func orEmptyString(data *string) string {
	if data != nil {
		return *data
	}
	return ""
}

// ToCSV serializes message to CSV row set.
// Statements are base for row set, rest of envelope data is duplicated in every row.
// Additional header row is added at the beginning.
func (m MT940Message) ToCSV(serializeT65 bool) []string {
	rows := []string{CSVHeader}
	for _, stmt := range m.Statements {
		row := []string{}
		row = append(row, m.TransactionRefNo)
		row = append(row, orEmptyString(m.RelatedReference))
		row = append(row, m.AccountIdentification.Account)
		row = append(row, orEmptyString(m.AccountIdentification.IdentCode))
		row = append(row, m.StatementNumber.StatementNo)
		row = append(row, orEmptyString(m.StatementNumber.SequenceNo))
		row = append(row,
			m.OpeningBalance.DCMark,
			m.OpeningBalance.Date.Format(time.DateOnly),
			m.OpeningBalance.Currency,
			m.OpeningBalance.Amount.StringFixed(2))
		row = append(row, stmt.Statement.ValueDate.Format(time.DateOnly))
		edate := ""
		if stmt.Statement.EntryDate != nil {
			edate = stmt.Statement.EntryDate.Format(time.DateOnly)
		}
		row = append(row, edate)
		row = append(row, stmt.Statement.DCMark)
		row = append(row, orEmptyString(stmt.Statement.FundsCode))
		row = append(row, stmt.Statement.Amount.StringFixed(2))
		row = append(row, stmt.Statement.TransactionIdent)
		row = append(row, stmt.Statement.Reference)
		row = append(row, orEmptyString(stmt.Statement.InstitutionReference))
		row = append(row, orEmptyString(stmt.Statement.Details))
		row = append(row, strings.Join(stmt.AccountOwnerInfo, " "))
		row = append(row,
			m.ClosingBalance.DCMark,
			m.ClosingBalance.Date.Format(time.DateOnly),
			m.ClosingBalance.Currency,
			m.ClosingBalance.Amount.StringFixed(2))
		if m.ClosingAvailableBalance != nil {
			row = append(row,
				m.ClosingAvailableBalance.DCMark,
				m.ClosingAvailableBalance.Date.Format(time.DateOnly),
				m.ClosingAvailableBalance.Currency,
				m.ClosingAvailableBalance.Amount.StringFixed(2))
		} else {
			row = append(row, "", "", "", "")
		}
		if serializeT65 {
			dc := []string{}
			dt := []string{}
			cur := []string{}
			amt := []string{}
			for _, fab := range m.ForwardAvailableBalance {
				dc = append(dc, fab.DCMark)
				dt = append(dt, fab.Date.Format(time.DateOnly))
				cur = append(cur, fab.Currency)
				amt = append(amt, fab.Amount.StringFixed(2))
			}
			row = append(row,
				strings.Join(dc, "/"),
				strings.Join(dt, "/"),
				strings.Join(cur, "/"),
				strings.Join(amt, "/"))
		} else {
			row = append(row, "", "", "", "")
		}
		row = append(row, strings.Join(m.AccountOwnerInfo, " "))
		rows = append(rows, strings.Join(row, ","))
	}
	return rows
}
