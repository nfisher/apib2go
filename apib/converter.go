package apib

import "github.com/nfisher/apib2go/primitives"

func String(s string) primitives.String {
	return primitives.String(&s)
}

func Number(n string) primitives.Number {
	return primitives.Number(&n)
}

func Boolean(b bool) primitives.Boolean {
	return primitives.Boolean(&b)
}
