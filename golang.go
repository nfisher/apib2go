package main

import (
	"fmt"
	"io"
	"strings"
)

type GoWriter struct {
	io.Writer
	open bool
	name string
}

type b []byte

func bs(format string, args ...interface{}) []byte {
	s := fmt.Sprintf(format, args...)
	return b(s)
}

func (gw *GoWriter) Append(item Item) {
	switch item.Type {
	case ItemModel:
		if gw.open {
			gw.Write(bs("}\n\n"))
		} else {
			gw.Write(bs("package %v\n\n", gw.name))
			gw.Write(bs("import . \"github.com/nfisher/apib2go/primitives\"\n\n"))
		}
		gw.Write(bs("type %s struct {\n", item.Value))
		gw.open = true
		return

	case ItemPropertyName:
		gw.Write(bs("  %v", strings.Title(item.Value)))
		return

	case ItemPropertyType:
		gw.Write(bs(" %v `json:\",omitempty\"`\n", strings.Title(item.Value)))
		return

	case ItemPropertyArrayType:
		gw.Write(bs(" []%v `json:\",omitempty\"`\n", strings.Title(item.Value)))
		return

	case ItemError:
		fmt.Printf("%s\n", item.Value)
		return
	}
}

func (gw *GoWriter) Close() {
	if gw.open {
		gw.Write(bs("}\n"))
	}
}
