package parser_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/oswida/mt9x/grammar"
	"github.com/oswida/mt9x/parser"
	"gotest.tools/v3/golden"
)

func TestProperMT940Files(t *testing.T) {
	parser := parser.NewFileParser[grammar.MT940Message]()
	basePath := filepath.Join("testdata", "mt940")
	files, err := os.ReadDir(filepath.Join(basePath, "input"))
	assert.NoError(t, err)
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		result, err := parser.Parse(filepath.Join(basePath, "input", f.Name()), false, nil)
		assert.NoError(t, err)
		value, err := json.MarshalIndent(result, "", " ")
		assert.NoError(t, err)
		golden.Assert(t, string(value), filepath.Join("mt940", "expected", strings.ReplaceAll(f.Name(), ".sta", ".json")))
	}

}
