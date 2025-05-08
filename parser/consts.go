package parser

const (
	// https://www2.swift.com/knowledgecentre/publications/usgi_20240719/2.0?topic=con_31519.htm
	StringX         = `[^:][a-zA-Z0-9/\-\?:().,'+  ]+`
	StringXNo2Slash = `[^:][a-zA-Z0-9\-\?:().,'+  ]+(/[a-zA-Z0-9\-\?:().,'+  ]+)*`
	Amount          = `[0-9]+,?[0-9]*` // this is MT9x amount with comma instead a dot
	CRLF            = "\r\n"
)
