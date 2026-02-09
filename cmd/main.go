package main

import (
	"fmt"
	"time"

	"github.com/Daemon0x00000000/lql/pkg/lql"
)

func main() {
	expr := &lql.Expression{}

	query := "!(sys_id=123^XORnameCONTAINS'example')"
	fmt.Printf("Query: %s\n\n", query)

	err := expr.Parse(query)
	if err != nil {
		panic(err)
	}

	data := map[string]interface{}{
		"sys_id": "123",
		"name":   "example",
	}

	since := time.Now()
	result, err := expr.Eval(data)
	if err != nil {
		panic(err)
	}
	elapsed := time.Since(since)

	fmt.Printf("Data: %v\n", data)
	fmt.Printf("Result: %v in %v\n", result, elapsed)
}
