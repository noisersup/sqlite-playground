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

func jsonbHandler(ctx *sqlite.FunctionContext, args []driver.Value) (driver.Value, error) {
	return base64.StdEncoding.EncodeToString([]byte(args[0].(string))), nil
}

func unmarshalHandler(ctx *sqlite.FunctionContext, args []driver.Value) (driver.Value, error) {
	return base64.StdEncoding.DecodeString(args[0].(string))
}
