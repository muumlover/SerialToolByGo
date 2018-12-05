package main

import "fmt"

type Parity string

const v Parity = "hello world"

func main() {
	fmt.Println(typeof(v))
}

func typeof(v interface{}) string {
	return fmt.Sprintf("%T", v)
}
