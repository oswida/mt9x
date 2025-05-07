package mt940

import (
	"mt9x/common"
)

// Grammar for MT940 file, according standard available here:
// https://www2.swift.com/knowledgecentre/publications/us9m_20230720/2.0?topic=mt940.htm

// Message represents MT940 standard message structure.
type Message struct {
	// Specifies the reference assigned by the Sender to unambiguously identify the message.
	TransactionRefNo string `parser:"T20 @StringX (CRLF|EOF)" json:"tag20"`
	// Contains the field 20 Transaction Reference Number of the request message (response to MT920 Request Message).
	RelatedReference *string `parser:"(T21 @StringX (CRLF|EOF))?" json:"tag21,omitempty"`
	// Identifies the account and optionally the identifier code of the account owner for which the statement is sent.
	// Need some examples, optional
	AccountIdentification AccountIdent `parser:"(T25|T25P) @@ (CRLF|EOF)" json:"tag25"`
	// Contains the sequential number of the statement, optionally followed by the sequence number of the message
	// within that statement when more than one message is sent for one statement.
	StatementNumber StatementNumber `parser:"T28C @@ (CRLF|EOF)" json:"tag28"`
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
	ForwardAvailableBalance *Balance `parser:"(T65 @@ (CRLF|EOF))?" json:"tag65,omitempty"`
	// Summarizing owner info
	AccountOwnerInfo []string `parser:"(T86 @StringX (CRLF @StringX)* (CRLF|EOF))?" json:"tag86,omitempty"`
}

type StatementSection struct {
	// Contains the details of each transaction.
	Statement Statement `parser:"T61 @@ (CRLF|EOF)" json:"tag61"`
	// Contains additional information about the transaction detailed in the preceding statement line
	// and which is to be passed on to the account owner.
	AccountOwnerInfo []string `parser:"(T86 @StringX (CRLF @StringX)* (CRLF|EOF))?" json:"tag86,omitempty"`
}

type AccountIdent struct {
	Account   string  `parser:"@StringX" json:"account"`
	IdentCode *string `parser:"(CRLF @StringX)?" json:"ident_code,omitempty"`
}

type StatementNumber struct {
	StatementNo string  `parser:"@Number" json:"stmt_number"`
	SequenceNo  *string `parser:"(Slash @Number)?" json:"seq_number,omitempty"`
}

type Balance struct {
	DCMark   string              `parser:"@DCMark" json:"dc_mark"`
	Date     common.SixDigitDate `parser:"@Date" json:"date"`
	Currency string              `parser:"@Currency" json:"currency"`
	Amount   common.CommaDecimal `parser:"@Amount" json:"amount"`
}

type Statement struct {
	ValueDate            common.SixDigitDate  `parser:"@Date" json:"value_date"`
	EntryDate            *common.SixDigitDate `parser:"@EntryDate?" json:"entry_date,omitempty"`
	DCMark               string               `parser:"@RDCMark" json:"dc_mark"`
	FundsCode            *string              `parser:"@BigLetter?" json:"funds_code,omitempty"`
	Amount               common.CommaDecimal  `parser:"@Amount" json:"amount"`
	TransactionIdent     string               `parser:"@TransIdent" json:"trx_ident"`
	Reference            string               `parser:"@StringXNoSlash" json:"owner_ref"`
	InstitutionReference *string              `parser:"(TwoSlashes @StringXNoSlash?)?" json:"institution_ref,omitempty"`
	Details              *string              `parser:"(CRLF @StringX)?" json:"details,omitempty"`
}
