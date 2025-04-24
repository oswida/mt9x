package mt940

// Grammar for MT940 file, according standard available here:
// https://www2.swift.com/knowledgecentre/publications/us9m_20230720/2.0?topic=mt940.htm

// Message represents MT940 standard message structure.
type Message struct {
	// Specifies the reference assigned by the Sender to unambiguously identify the message.
	TransactionRefNo string `parser:"T20 @StringX (CRLF|EOF)"`
	// Contains the field 20 Transaction Reference Number of the request message (response to MT920 Request Message).
	RelatedReference *string `parser:"(T21 @StringX (CRLF|EOF))?"`
	// Identifies the account and optionally the identifier code of the account owner for which the statement is sent.
	// Need some examples, optional
	AccountIdentification AccountIdent `parser:"(T25|T25P) @@ (CRLF|EOF)"`
	// Contains the sequential number of the statement, optionally followed by the sequence number of the message
	// within that statement when more than one message is sent for one statement.
	StatementNumber StatementNumber `parser:"T28C @@ (CRLF|EOF)"`
	// Specifies, for the (intermediate - M) opening balance, whether it is a debit or credit balance,
	// the date, the currency and the amount of the balance.
	OpeningBalance Balance `parser:"(T60F|T60M) @@ (CRLF|EOF)"`
	// Statement information
	Statements []StatementSection `parser:"@@*"`
	// Specifies, for the (intermediate) closing balance.
	ClosingBalance Balance `parser:"(T62F|T62M) @@ (CRLF|EOF)"`
	// Indicates the funds which are available to the account owner (if credit balance)
	// or the balance which is subject to interest charges (if debit balance).
	ClosingAvailableBalance *Balance `parser:"(T64 @@ (CRLF|EOF))?"`
	// Indicates the funds which are available to the account owner
	// (if a credit or debit balance) for the specified forward value date.
	ForwardAvailableBalance *Balance `parser:"(T65 @@ (CRLF|EOF))?"`
	// Summarizing owner info
	AccountOwnerInfo []string `parser:"(T86 @StringX (CRLF @StringX)* (CRLF|EOF))?"`
}

type StatementSection struct {
	// Contains the details of each transaction.
	Statement Statement `parser:"T61 @@ (CRLF|EOF)"`
	// Contains additional information about the transaction detailed in the preceding statement line
	// and which is to be passed on to the account owner.
	AccountOwnerInfo []string `parser:"(T86 @StringX (CRLF @StringX)* (CRLF|EOF))?"`
}

type AccountIdent struct {
	Account   string  `parser:"@StringX"`
	IdentCode *string `parser:"(CRLF @StringX)?"`
}

type StatementNumber struct {
	StatementNo string  `parser:"@Number"`
	SequenceNo  *string `parser:"(Slash @Number)?"`
}

type Balance struct {
	DCMark   string `parser:"@DCMark"`
	Date     string `parser:"@Date"`
	Currency string `parser:"@Currency"`
	Amount   string `parser:"@Amount"`
}

type Statement struct {
	ValueDate            string  `parser:"@Date"`
	EntryDate            *string `parser:"@EntryDate?"`
	DCMark               string  `parser:"@RDCMark"`
	FundsCode            *string `parser:"@BigLetter?"`
	Amount               string  `parser:"@Amount"`
	TransactionIdent     string  `parser:"@TransIdent"`
	Reference            string  `parser:"@StringXNoSlash"`
	InstitutionReference *string `parser:"(TwoSlashes @StringXNoSlash?)?"`
	Details              *string `parser:"(CRLF @StringX)?"`
}
