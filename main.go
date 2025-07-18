package main

import (
	"encoding/json"
	"fmt"

	"github.com/oswida/mt9x/grammar"
	"github.com/oswida/mt9x/parser"
)

func main() {
	p := parser.NewFileParser[grammar.MT940Message]()
	result, err := p.Parse("parser/testdata/mt940/input/csob.sta", false, nil)
	if err != nil {
		panic(err)
	}
	v, err := json.MarshalIndent(result, "", " ")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", string(v))
}
