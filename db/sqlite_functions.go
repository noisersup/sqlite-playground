package db

import (
	"database/sql/driver"
	"encoding/base64"

	"modernc.org/sqlite"
)

// INSERT INTO test VALUES(jsonb({"v":1}))
// SELECT unmarshal(x) FROM test
func initSQLiteFunctions() {
	sqlite.MustRegisterScalarFunction("jsonb2", 1, jsonbHandler)
	sqlite.MustRegisterScalarFunction("unmarshal", 1, unmarshalHandler)
}

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
	return base64.StdEncoding.EncodeToString([]byte(args[0].(string))), nil
}

func unmarshalHandler(ctx *sqlite.FunctionContext, args []driver.Value) (driver.Value, error) {
	return base64.StdEncoding.DecodeString(args[0].(string))
}

type jsonType uint8

const (
	keyTyp jsonType = iota
	stringTyp
	numTyp
	floatTyp
)
