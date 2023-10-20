package db

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/Jeffail/gabs/v2"
	"modernc.org/sqlite"
)

var (
	contType      jsonType = 0b000
	keyType       jsonType = 0b001
	stringType    jsonType = 0b010
	numericType   jsonType = 0b011
	boolTrueType  jsonType = 0b100
	boolFalseType jsonType = 0b101
	nullType      jsonType = 0b110
	containerType jsonType = 0b111
)

// NOTE: it seems impossible to implement jsonb in the same way as postgres does.
// While pg's jsonb is able to store data in a tree structure where every node is varlena,
// SQLite doesn't give any tools to introduce own data types so the only option is to use existing types
// and marshal them via registered functions.
// The main downside is the fact that every document needs to be placed in memory before we start do do
// unmarshaling.
//
// Aside from that there's still a lot of room to improve performence by storing the json types and bitlength of every item
// in "headers" that might take ~/<byte.
// By doing so our helper jsonb functions wouldn't need to marshal whole json every time, but rather just "jump" from header to header
//
// This potentially might also introduce some way of indexing specific keys in documents by using expression indexes and virtual columns(?)
func jsonbHandler(ctx *sqlite.FunctionContext, args []driver.Value) (driver.Value, error) {
	input := args[0].(string)

	m := map[string]interface{}{}

	if err := json.Unmarshal([]byte(input), &m); err != nil {
		return nil, err
	}

	var res []byte

	for k, v := range m {
		// 11100000 - dataType mask (0-7)
		// 00011111 - dataLen in bytes mask (0-31)
		length := len(k)
		res = append(res, genHeader(keyType, uint8(length)))
		res = append(res, marshalValue(k)...)

		val := marshalValue(v)
		res = append(res, genHeader(getType(v), uint8(len(val))))
		res = append(res, val...)
	}

	return res, nil
}

func marshalValue(v any) []byte {
	switch getType(v) {
	case stringType:
		return []byte(v.(string))
	default:
		panic("not yet, how did you access it tho?")
	}
}

func getType(v any) jsonType {
	switch v.(type) {
	case string:
		return stringType
	default:
		panic("not yet")
	}
}

// TODO make iterator here and make it work overall when you're less tired lol
func unmarshalHandler(ctx *sqlite.FunctionContext, args []driver.Value) (driver.Value, error) {
	input, ok := args[0].([]byte)
	if !ok {
		return nil, fmt.Errorf("Invalid type, expected []byte")
	}

	jsonObj := gabs.New()
	iter := NewIterator(input)

	for {
		b, err := iter.Next()
		if err == ErrIteratorDone {
			break
		}
		if err != nil {
			return nil, err
		}

		typ, length := readHeader(b)
		if typ != keyType {
			panic("huh")
		}

		var key []byte
		for i := 0; i < int(length); i++ {
			keyByte, err := iter.Next()
			if err != nil {
				return nil, err
			}

			key = append(key, keyByte)
		}

		b, err = iter.Next()
		if err != nil {
			return nil, err
		}

		typ, length = readHeader(b)

		var v []byte
		for i := 0; i < int(length); i++ {
			valByte, err := iter.Next()
			if err != nil {
				return nil, err
			}

			v = append(v, valByte)
		}

		jsonObj.Set(string(v), string(key))
	}

	return jsonObj.String(), nil
}

func readHeader(header byte) (jsonType, uint8) {
	return jsonType(header & 0b111),
		header & 0b00011111
}

func genHeader(typ jsonType, length uint8) (header byte) {
	header = byte(typ << 5)
	if length > 31 || length < 0 {
		panic("Invalid length")
	}

	header += length
	return
}

type jsonType uint8
