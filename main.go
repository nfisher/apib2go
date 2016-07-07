package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	var filename string
	var pkgname string
	flag.StringVar(&filename, "input", "", "Input filename.")
	flag.StringVar(&pkgname, "package", "", "Package name.")
	flag.Parse()

	if filename == "" || pkgname == "" {
		flag.Usage()
		os.Exit(1)
	}

	r, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	l := New(filename, string(b))

	go func() {
		l.Run()
	}()

	w := &GoWriter{
		os.Stdout,
		false,
		pkgname,
		"",
	}
	defer w.Close()

	for item := range l.Items {
		// this is dirty should compose Items into AST nodes.
		w.Append(item)
	}
}
