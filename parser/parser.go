package parser

import (
	"fmt"
	"io"
	"os"

	"github.com/alecthomas/participle/v2"
)

type FileParser[T MT9xMessage] struct {
	parser *participle.Parser[T]
}

// NewFileParser creates new file parser for MT940 messages.
func NewFileParser[T MT9xMessage]() *FileParser[T] {
	lexer := NewLexer()
	parser := participle.MustBuild[T](
		participle.Lexer(lexer),
		participle.UseLookahead(2))
	return &FileParser[T]{
		parser: parser,
	}
}

// Parse parses MT940 message into structure.
func (fp *FileParser[T]) Parse(filename string, validate bool, traceWriter io.Writer) (*T, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filename, err)
	}
	options := []participle.ParseOption{participle.AllowTrailing(true)}
	if traceWriter != nil {
		options = append(options, participle.Trace(traceWriter))
	}
	res, err := fp.parser.Parse(filename, f, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file %s: %w", filename, err)
	}
	if validate {
		if err = (*res).Validate(); err != nil {
			return nil, fmt.Errorf("failed to validate parsed result %s: %w", filename, err)
		}
	}

	return res, nil
}
