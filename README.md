# apib2go

[![Build Status](https://travis-ci.org/nfisher/apib2go.svg?branch=master)](https://travis-ci.org/nfisher/apib2go)
[![Coverage Status](https://coveralls.io/repos/github/nfisher/apib2go/badge.svg?branch=master)](https://coveralls.io/github/nfisher/apib2go?branch=master)

## Core Principles

- Optional is represented by pointer.
- Required is represented by an instance. \*
- Numbers are represented as strings to avoid issues with precision.
- functions such as builders and other tooling can co-exist with the generated code as long as the user does not place it in the generated source code files.

\* I'm not really a fan of required fields. It limits schema evolution (protobuf3 dropped it entirely) so it'll be the last thing I focus on.

## Type Mapping

| APIB    | Go            | Notes                               |
| ------- | ------------- | ----------------------------------- |
| boolean | \*bool        |                                     |
| string  | \*string      |                                     |
| number  | \*string      | allow user to decide how to convert |
| array   | pointer slice |                                     |
| enum    | \*string      |                                     |
| object  | \*Object      |                                     |

## Example

fruits.apib
```
# Fruit API
Fruit distribution API.

## Data Structures

### Dimension
+ radius (number)
+ length (number)

### Produce

+ colour (string) - What colour is it?
+ dimesions (Dimension)
+ fruit (boolean) - Is it fruit?
```

Execution
```
apib2go -input fruit.apib -package fruit
```


Go
```
package fruit

import (
  "github.com/nfisher/apib2go/apib"
  . "github.com/nfisher/apib2go/primitives"
)

type Dimensions struct {
  Radius Number
  Length Number
}

type Produce struct {
  Colour     String
  Dimensions *Dimensions
  Fruit      Boolean
  Name       String
}

// Usage:

dim := &Dimension {
  Radius: apib.Int(3),
  Length: apib.Decimal("18.56"),
}

p := &Produce {
  Colour: apib.String("yellow"),
  Dimensions: dim,
  Fruit: apib.Boolean(true),
  Name: apib.String("banana"),
}
```

## Reference Material

- https://apiblueprint.org/documentation/specification.html
- https://apiblueprint.org/documentation/mson/specification.html
