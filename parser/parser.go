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

type ByteParser[T MT9xMessage] struct {
	parser *participle.Parser[T]
}

// NewByteParser creates new byte parser for MT940 messages.
func NewByteParser[T MT9xMessage]() *ByteParser[T] {
	lexer := NewLexer()
	parser := participle.MustBuild[T](
		participle.Lexer(lexer),
		participle.UseLookahead(2))
	return &ByteParser[T]{
		parser: parser,
	}
}

// Parse parses MT940 message (from data) into structure.
func (fp *ByteParser[T]) Parse(data []byte, validate bool, traceWriter io.Writer) (*T, error) {
	options := []participle.ParseOption{participle.AllowTrailing(true)}
	if traceWriter != nil {
		options = append(options, participle.Trace(traceWriter))
	}
	res, err := fp.parser.ParseBytes("byte data", data, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to parse bytes: %w", err)
	}
	if validate {
		if err = (*res).Validate(); err != nil {
			return nil, fmt.Errorf("failed to validate parsed result: %w", err)
		}
	}

	return res, nil
}
