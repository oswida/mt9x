package main

import (
	"encoding/json"
	"fmt"
	"mt9x/mt940"
)

func main() {
	p := mt940.NewFileParser()
	result, err := p.Parse("testdata/msg/spec-example.msg")
	if err != nil {
		panic(err)
	}
	v, err := json.MarshalIndent(result, "", " ")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", string(v))
}
