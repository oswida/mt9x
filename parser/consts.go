package parser

const (
	StringX        = `[^:][a-zA-Z0-9/\-\?:().,'+  ]+`
	StringXNoSlash = `[^:][a-zA-Z0-9\-\?:().,'+  ]+`
	Amount         = `[0-9]+,?[0-9]*` // this is MT9x amount with comma instead a dot
	CRLF           = "\r\n"
)
