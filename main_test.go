package main

import (
	"flag"
	"log"
	"sqlite-battlefield/db"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite" // register database/sql driver
)

var jsonbF = flag.Bool("jsonb", false, "run tests against jsonb backend")

// TestMain is the entry point for all integration tests.
func TestMain(m *testing.M) {
	flag.Parse()

	m.Run()
}

//func init() {
//	sqlite.MustRegisterScalarFunction("jsonb2", 1, jsonbHandler)
//	sqlite.MustRegisterScalarFunction("unmarshal", 1, unmarshalHandler)
//}

func GetTestDB() db.DB {
	if *jsonbF {
		return db.NewJsonbDB()
	}
	return db.NewJsonDB()
}

func TestJson(t *testing.T) {
	db := GetTestDB()

	for _, jsonStr := range []string{
		//`{}`,
		`{"v":"foo"}`,
		//`{"v":42}`,
		//`{"v":"42"}`,
		//`{"v":42.13}`,
		//`{"v":"42.13"}`,
		//`{"_id":1,"v":42}`,
	} {
		require.NoError(t, db.DropTable())
		require.NoError(t, db.CreateTable())

		require.NoError(t, db.Insert(jsonStr))

		res, err := db.Query()
		require.NoError(t, err)

		log.Println(res)
		assert.Equal(t, jsonStr, res)
	}
}
