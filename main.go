package main

import (
	"encoding/json"
	"fmt"
	"mt9x/grammar"
	"mt9x/parser"
)

func main() {
	p := parser.NewFileParser[grammar.MT940Message]()
	result, err := p.Parse("testdata/mt940/input/spec-example.msg", nil)
	if err != nil {
		panic(err)
	}
	v, err := json.MarshalIndent(result, "", " ")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", string(v))
}
