package main

import (
	"flag"
	"fmt"
)

func main() {
	var filename string
	flag.StringVar(&filename, "input", "", "Input filename.")
	flag.Parse()
	fmt.Println(filename)
}
