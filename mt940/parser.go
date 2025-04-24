package mt940

import (
	"fmt"
	"os"

	"github.com/alecthomas/participle/v2"
)

type FileParser struct {
	parser *participle.Parser[Message]
}

// NewFileParser creates new file parser for MT940 messages.
func NewFileParser() *FileParser {
	lexer := NewLexer()
	parser := participle.MustBuild[Message](
		participle.Lexer(lexer),
		participle.UseLookahead(1))
	return &FileParser{
		parser: parser,
	}
}

// Parse parses MT940 message into structure.
func (fp *FileParser) Parse(filename string) (*Message, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filename, err)
	}
	result, err := fp.parser.Parse(filename, f)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file %s: %w", filename, err)
	}
	err = result.Validate()
	if err != nil {
		return nil, fmt.Errorf("failed to validate parsed result %s: %w", filename, err)
	}

	return result, nil
}
