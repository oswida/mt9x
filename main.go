package main

import (
	"encoding/json"
	"fmt"
	"mt9x/grammar"
	"mt9x/parser"
	"os"
)

func main() {
	p := parser.NewFileParser[grammar.MT940Message]()
	result, err := p.Parse("parser/testdata/mt940/input/ok-mbank.msg", os.Stdout)
	if err != nil {
		panic(err)
	}
	v, err := json.MarshalIndent(result, "", " ")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", string(v))
}
