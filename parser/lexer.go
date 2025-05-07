package parser

import (
	"github.com/alecthomas/participle/v2/lexer"
)

// NewLexer creates stateful lexical analyzer for MT9x messages.
func NewLexer() *lexer.StatefulDefinition {
	return lexer.MustStateful(lexer.Rules{
		"Root": []lexer.Rule{
			{Name: "Slash", Pattern: "/", Action: nil},
			{Name: "CRLF", Pattern: CRLF, Action: nil},
			{Name: "T20", Pattern: ":20:", Action: lexer.Push("Tag20")},
			{Name: "T21", Pattern: ":21:", Action: lexer.Push("Tag21")},
			{Name: "T25", Pattern: ":25:", Action: lexer.Push("Tag25")},
			{Name: "T25P", Pattern: ":25P:", Action: lexer.Push("Tag25")},
			{Name: "T28C", Pattern: ":28C:", Action: lexer.Push("Tag28")},
			{Name: "T60F", Pattern: ":60F:", Action: lexer.Push("Tag60")},
			{Name: "T60M", Pattern: ":60M:", Action: lexer.Push("Tag60")},
			{Name: "T61", Pattern: ":61:", Action: lexer.Push("Tag61")},
			{Name: "T86", Pattern: ":86:", Action: lexer.Push("Tag86")},
			{Name: "T62F", Pattern: ":62F:", Action: lexer.Push("Tag60")},
			{Name: "T62M", Pattern: ":62M:", Action: lexer.Push("Tag60")},
			{Name: "T64", Pattern: ":64:", Action: lexer.Push("Tag60")},
			{Name: "T65", Pattern: ":65:", Action: lexer.Push("Tag60")},
		},
		"Tag20": []lexer.Rule{
			{Name: "StringX", Pattern: StringX, Action: nil},
			{Name: "CRLF", Pattern: CRLF, Action: lexer.Pop()},
		},
		"Tag21": []lexer.Rule{
			{Name: "StringX", Pattern: StringX, Action: nil},
			{Name: "CRLF", Pattern: CRLF, Action: lexer.Pop()},
		},
		"Tag25": []lexer.Rule{
			{Name: "StringX", Pattern: StringX, Action: nil},
			// Possible identifier
			{Name: "CRLF", Pattern: CRLF, Action: lexer.Push("Tag25_1")},
			lexer.Return(),
		},
		"Tag25_1": []lexer.Rule{
			{Name: "StringX", Pattern: StringX, Action: nil}, // cannot be a tag
			{Name: "CRLF", Pattern: CRLF, Action: lexer.Pop()},
			lexer.Return(),
		},
		"Tag28": []lexer.Rule{
			{Name: "Number", Pattern: "[0-9]+", Action: nil},
			{Name: "Slash", Pattern: "/", Action: nil},
			{Name: "CRLF", Pattern: CRLF, Action: lexer.Pop()},
		},
		"Tag60": []lexer.Rule{
			{Name: "DCMark", Pattern: "[DC]", Action: nil},
			{Name: "Currency", Pattern: "[A-Z][A-Z][A-Z]", Action: lexer.Push("Tag60Amount")},
			{Name: "Date", Pattern: "[0-9][0-9][0-9][0-9][0-9][0-9]", Action: nil},
			{Name: "CRLF", Pattern: CRLF, Action: lexer.Pop()},
		},
		"Tag60Amount": []lexer.Rule{
			// Less strict as in standard as it accepts also amounts w/o comma
			{Name: "Amount", Pattern: Amount, Action: lexer.Pop()},
			{Name: "CRLF", Pattern: CRLF, Action: lexer.Pop()},
		},
		"Tag61": []lexer.Rule{
			// recognize first, mandatory date
			{Name: "Date", Pattern: "[0-9][0-9][0-9][0-9][0-9][0-9]", Action: lexer.Push("Tag61_1")},
			// (Reversed) Debit/Credit mark
			{Name: "RDCMark", Pattern: "[R]?[DC]", Action: lexer.Push("Tag61_2")},
			// Maybe end but also supplementary details
			{Name: "CRLF", Pattern: CRLF, Action: nil},
			lexer.Return(),
		},
		"Tag61_1": []lexer.Rule{
			// Optional entry date
			{Name: "EntryDate", Pattern: "[0-9][0-9][0-9][0-9]", Action: nil},
			lexer.Return(),
		},
		"Tag61_2": []lexer.Rule{
			// Optional funds code
			{Name: "BigLetter", Pattern: "[A-Z]", Action: nil},
			// Transaction amount
			// Less strict as in standard as it accepts also amounts w/o comma
			{Name: "Amount", Pattern: Amount, Action: lexer.Push("Tag61_3")},
			lexer.Return(),
		},
		"Tag61_3": []lexer.Rule{
			// Transaction identification
			{Name: "TransIdent", Pattern: "[SNF][A-Z0-9][A-Z0-9][A-Z0-9]", Action: lexer.Push("Tag61_4")},
			lexer.Return(),
		},
		"Tag61_4": []lexer.Rule{
			// References
			{Name: "StringXNoSlash", Pattern: StringXNoSlash, Action: nil},
			{Name: "TwoSlashes", Pattern: "//", Action: nil},
			{Name: "CRLF", Pattern: CRLF, Action: lexer.Push("Tag61_5")},
			lexer.Return(),
		},
		"Tag61_5": []lexer.Rule{
			// Supplementary details
			{Name: "StringX", Pattern: StringX, Action: nil}, // cannot be tag
			lexer.Return(),
		},
		"Tag86": []lexer.Rule{
			// this tag can have 6 lines maximum
			{Name: "StringX", Pattern: StringX, Action: nil}, // cannot be tag
			{Name: "CRLF", Pattern: CRLF, Action: nil},
			lexer.Return(),
		},
	})
}
