package mt940

import (
	"fmt"
	"strings"
)

func (m *Message) String() string {
	result := []string{}
	result = append(result, "transaction ref no: "+m.TransactionRefNo)
	if m.RelatedReference != nil {
		result = append(result, "related reference: "+*m.RelatedReference)
	}
	ident := m.AccountIdentification.Account
	if m.AccountIdentification.IdentCode != nil {
		ident = ident + ":" + *m.AccountIdentification.IdentCode
	}
	result = append(result, "account identification: "+ident)
	result = append(result, "statement number: "+m.StatementNumber.String())
	result = append(result, "opening balance: "+m.OpeningBalance.String())

	for _, st := range m.Statements {
		result = append(result, "statement: "+st.String())
	}

	result = append(result, "closing balance: "+m.ClosingBalance.String())
	if m.ClosingAvailableBalance != nil {
		result = append(result, "closing available balance: "+m.ClosingAvailableBalance.String())
	}
	if m.ForwardAvailableBalance != nil {
		result = append(result, "forwarded available balance: "+m.ForwardAvailableBalance.String())
	}
	if m.AccountOwnerInfo != nil {
		result = append(result, "account owner info: "+strings.Join(m.AccountOwnerInfo, " "))
	}

	return strings.Join(result, "\n")
}

func (sn StatementNumber) String() string {
	result := sn.StatementNo
	if sn.SequenceNo != nil {
		result = result + "/" + *sn.SequenceNo
	}

	return result
}

func (b *Balance) String() string {
	result := b.DCMark + " " + b.Date + " " + b.Amount + " " + b.Currency

	return result
}

func (ss *StatementSection) String() string {
	result := ss.Statement.String()
	if ss.AccountOwnerInfo != nil {
		result = result + "[" + strings.Join(ss.AccountOwnerInfo, ", ") + "]"
	}

	return result
}

func (s *Statement) String() string {
	entryDate := ""
	if s.EntryDate != nil {
		entryDate = *s.EntryDate
	}
	fcode := ""
	if s.FundsCode != nil {
		fcode = *s.FundsCode
	}
	iref := ""
	if s.InstitutionReference != nil {
		iref = *s.InstitutionReference
	}
	details := ""
	if s.Details != nil {
		details = *s.Details
	}
	result := fmt.Sprintf("value date: %s; entry date: %s; DC mark: %s; funds code: %s; amount: %s; ident: %s; ref: %s; institution ref: %s; details: %s",
		s.ValueDate, entryDate, s.DCMark, fcode, s.Amount, s.TransactionIdent, s.Reference, iref, details)

	return result
}
