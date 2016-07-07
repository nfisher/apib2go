package main

import (
	"fmt"
	"io"
	"strings"
)

type GoWriter struct {
	io.Writer
	pkgname string
}

type b []byte

func bs(format string, args ...interface{}) []byte {
	s := fmt.Sprintf(format, args...)
	return b(s)
}

func (w *GoWriter) WriteDoc(doc *Document) {
	w.Write(bs("package %v\n\n", w.pkgname))
	w.Write(bs("import . \"github.com/nfisher/apib2go/primitives\"\n\n"))

	for _, model := range doc.DataStructures {
		w.Write(bs("type %s struct {\n", model.Name))
		for _, property := range model.Properties {
			pre := "*"
			if strings.Contains("string number boolean ", property.Type) {
				pre = ""
			}
			f := "  %v %v%v `json:\"%v,omitempty\"`\n"
			if property.IsArray {
				f = "  %v []%v%v `json:\"%v,omitempty\"`\n"
			}
			w.Write(bs(f, strings.Title(property.Name), pre, strings.Title(property.Type), property.Name))
		}
		w.Write(bs("}\n\n"))
	}
}
