package parser

// https://www2.swift.com/knowledgecentre/publications/usgi_20240719/2.0?topic=con_31492.htm

const (
	// Actually we are using Z character set instead of X, some banks are not strict to the char sets
	CharX        = `[a-zA-Z0-9/\-\?:().,'+  ="%&\*<>;@#_]`
	CharXNoSlash = `[a-zA-Z0-9\-\?:().,'+  ="%&\*<>;@#_]`
	Numeric      = `[0-9]`               // n
	AlphaUpper   = `[A-Z]`               // a
	AlphaNum     = `[A-Z0-9]`            // c
	HexNum       = `[ABCDEF0-9]`         // h
	CharXSeq     = `^[^:]` + CharX + `+` // x sequence
	// This field must not start or end with a slash '/' and must not contain two consecutive slashes '//' (Error code(s): T26).
	CharXSeqSlashRestrict = `^[^:]` + CharXNoSlash + `+(/` + CharXNoSlash + `+)*`
	NumSeq                = Numeric + `+`
	DCMark                = `[DC]`
	RDCMark               = `R?[DC]`
	Numeric46             = Numeric + Numeric + Numeric + `(?:` + Numeric + Numeric + Numeric + `|` + Numeric + `)`
	Alpha3                = AlphaUpper + AlphaUpper + AlphaUpper
	Amount                = Numeric + `+(,` + Numeric + `*)?` // this is MT9x amount with comma instead a dot
	TrxIdentCode          = `(?:S` + Numeric + Numeric + Numeric + `|[NF]` + AlphaNum + AlphaNum + AlphaNum + `)`
	CRLF                  = "\r\n"
)
