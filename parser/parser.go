package parser

import (
	"fmt"
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
		participle.UseLookahead(1))
	return &FileParser[T]{
		parser: parser,
	}
}

// Parse parses MT940 message into structure.
func (fp *FileParser[T]) Parse(filename string) (*T, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filename, err)
	}
	res, err := fp.parser.Parse(filename, f)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file %s: %w", filename, err)
	}
	err = (*res).Validate()
	if err != nil {
		return nil, fmt.Errorf("failed to validate parsed result %s: %w", filename, err)
	}

	return res, nil
}
