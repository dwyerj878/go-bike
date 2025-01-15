package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	println("Hello")
	println(os.Args[1])
	fileName := os.Args[1]

	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	var result map[string]interface{}
	decoder := json.NewDecoder(file)

	err = decoder.Decode(&result)
	if err != nil {
		panic(err)
	}
	fmt.Print(result["RIDE"])
}
