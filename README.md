# Go packed serializer [![Build Status](https://travis-ci.org/ikkerens/serialize.svg?branch=master)](https://travis-ci.org/ikkerens/serialize) [![Go Report Card](https://goreportcard.com/badge/github.com/ikkerens/serialize)](https://goreportcard.com/report/github.com/ikkerens/serialize) [![GoDoc](https://godoc.org/github.com/ikkerens/serialize?status.svg)](https://godoc.org/github.com/ikkerens/serialize)

This is a packed struct serializer that is mostly meant for a private project but was released as it may be useful to someone else.

Originally this package was made as an extension to binary.Read and binary.Write, but I soon found those functions didn't match my use case as they offered no support for strings nor compression.

#### Features
* Caches types for faster calls to the same type
* Compression support
* Tread safe (the calls are, reading to the value is not)
* Easy to implement in other languages
* Supported types:
  * uint, uint8 (and byte) up to uint64
  * int, int8 up to int64
  * float32 and float64
  * string
  * anything implementing the Serializer/Deserializer interfaces
  * slices
  * structs

#### Format
* All primitives are stored in big endian format
* All slices are stored with a uint32 prefix indicating their length
* Strings are stored with a uint32 prefix indicating their length
* Compression blocks are stored using deflate (level 9) with a uint32 prefixing the size of the compressed data blob

#### Note about int/uint
The types `int` and `uint` will be treated as `int64` and `uint64` respectively, as their actual sizes are
dependant on the compiler architecture used. However, they are not recommended as the receiving end may not fit
storing 64 bits in their int type, thus overflowing.

## Include in your project
```go
import "github.com/ikkerens/serialize"
```

## Usage
```go
package main

import (
	"bytes"
	"log"

	"github.com/ikkerens/serialize"
)

type myBlob struct {
	A uint64 // all fields have to be exported
	B []byte `compressed:"true"` // this field will be serialized compressed, can be added anywhere
	C subBlob
}

type subBlob struct {
	D string
}

func main() {
	b := new(bytes.Buffer)
	blob := &myBlob{A: 1, B: []byte{1, 2, 3, 4}, C: subBlob{D: "test message"}}

	// Serialize
	if err := serialize.Write(b, blob); err != nil { // Write does not need a pointer, but it is recommended
		log.Fatalln(err)
	}

	// Deserialize
	newBlob := new(myBlob)
	if err := serialize.Read(b, newBlob); err != nil { // Read *needs* a pointer, or it will panic
		log.Fatalln(err)
	}

	log.Printf("Successfully deserialized: %+v", newBlob)
}
```
