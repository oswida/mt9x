package parser

import (
	"github.com/alecthomas/participle/v2/lexer"
)

// NewLexer creates stateful lexical analyzer for MT9x messages.
func NewLexer() *lexer.StatefulDefinition {
	return lexer.MustStateful(lexer.Rules{
		"Root": []lexer.Rule{
			{Name: "T20", Pattern: ":20:", Action: lexer.Push("SlashRestricted")},
			{Name: "T21", Pattern: ":21:", Action: lexer.Push("SlashRestricted")},
			{Name: "T25", Pattern: ":25:", Action: lexer.Push("OnlyChars")},
			{Name: "T25P", Pattern: ":25P:", Action: lexer.Push("OnlyChars")},
			{Name: "T28C", Pattern: ":28C:", Action: lexer.Push("StmtNumber")},
			{Name: "T60F", Pattern: ":60F:", Action: lexer.Push("Balance_1")},
			{Name: "T60M", Pattern: ":60M:", Action: lexer.Push("Balance_1")},
			{Name: "T62F", Pattern: ":62F:", Action: lexer.Push("Balance_1")},
			{Name: "T62M", Pattern: ":62M:", Action: lexer.Push("Balance_1")},
			{Name: "T64", Pattern: ":64:", Action: lexer.Push("Balance_1")},
			{Name: "T65", Pattern: ":65:", Action: lexer.Push("Balance_1")},
			{Name: "T61", Pattern: ":61:", Action: lexer.Push("Statement_1")},
			{Name: "T86", Pattern: ":86:", Action: nil},
			{Name: "CRLF", Pattern: CRLF, Action: nil},
			{Name: "CharXSeq", Pattern: CharXSeq, Action: nil},
		},
		"SlashRestricted": []lexer.Rule{
			{Name: "CharXSeqSlashRestrict", Pattern: CharXSeqSlashRestrict, Action: nil},
			{Name: "CRLF", Pattern: CRLF, Action: nil},
			lexer.Return(),
		},
		"OnlyChars": []lexer.Rule{
			{Name: "CharXSeq", Pattern: CharXSeq, Action: nil},
			{Name: "CRLF", Pattern: CRLF, Action: nil},
			lexer.Return(),
		},
		"StmtNumber": {
			{Name: "Slash", Pattern: "/", Action: nil},
			{Name: "NumSeq", Pattern: NumSeq, Action: nil},
			lexer.Return(),
		},
		"Balance_1": []lexer.Rule{
			{Name: "DCMark", Pattern: DCMark, Action: nil},
			{Name: "Date", Pattern: Numeric46, Action: nil},
			{Name: "Currency", Pattern: Alpha3, Action: lexer.Push("Balance_2")},
			lexer.Return(),
		},
		"Balance_2": {
			{Name: "Amount", Pattern: Amount, Action: lexer.Pop()},
			lexer.Return(),
		},
		"Statement_1": []lexer.Rule{
			{Name: "Date", Pattern: Numeric46, Action: nil}, // check with EntryDate
			{Name: "RDCMark", Pattern: RDCMark, Action: lexer.Push("Statement_2")},
			{Name: "BigLetter", Pattern: AlphaUpper, Action: lexer.Push("Statement_2")},
			lexer.Return(),
		},
		"Statement_2": []lexer.Rule{
			{Name: "Amount", Pattern: Amount, Action: nil},
			{Name: "TransIdent", Pattern: TrxIdentCode, Action: lexer.Push("Statement_3")},
			lexer.Return(),
		},
		"Statement_3": []lexer.Rule{
			{Name: "CharXSeqSlashRestrict", Pattern: CharXSeqSlashRestrict, Action: nil},
			{Name: "TwoSlashes", Pattern: "//", Action: nil},
			{Name: "CRLF", Pattern: CRLF, Action: lexer.Push("OnlyChars")},
			lexer.Return(),
		},
	})
}
