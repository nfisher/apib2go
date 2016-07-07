package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

type Property struct {
	Description string
	Name        string
	Type        string
	IsArray     bool
}

type DataStructure struct {
	Name       string
	Properties []*Property
}

type MetaData struct {
	Key   string
	Value string
}

type Document struct {
	MetaData       []*MetaData
	DataStructures []*DataStructure
}

func NewDoc() *Document {
	return &Document{
		MetaData:       make([]*MetaData, 0, 10),
		DataStructures: make([]*DataStructure, 0, 10),
	}
}

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

	doc := NewDoc()

	var md *MetaData
	var model *DataStructure
	var prop *Property

	for item := range l.Items {
		switch item.Type {
		case ItemError:
			fmt.Printf("%v\n", item.Value)
			os.Exit(1)

		case ItemMetaKey:
			md = &MetaData{}
			md.Key = item.Value
			continue

		case ItemMetaValue:
			md.Value = item.Value
			doc.MetaData = append(doc.MetaData, md)
			md = nil
			continue

		case ItemModel:
			model = &DataStructure{}
			doc.DataStructures = append(doc.DataStructures, model)
			model.Name = item.Value
			continue

		case ItemPropertyName:
			prop = &Property{}
			model.Properties = append(model.Properties, prop)
			prop.Name = item.Value
			continue

		case ItemPropertyType:
			prop.Type = item.Value
			prop.IsArray = false
			continue

		case ItemPropertyArrayType:
			prop.Type = item.Value
			prop.IsArray = true
			continue
		}
	}

	w := &GoWriter{
		os.Stdout,
		pkgname,
	}

	w.WriteDoc(doc)
}
