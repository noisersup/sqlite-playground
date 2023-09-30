package db

import (
	"database/sql"
	"log"

	"modernc.org/sqlite"
	_ "modernc.org/sqlite" // register database/sql driver
)

type jsonbDB struct {
	db *sql.DB
}

func init() {
	sqlite.MustRegisterScalarFunction("jsonb", 1, jsonbHandler)
	sqlite.MustRegisterScalarFunction("unmarshal", 1, unmarshalHandler)
}

func NewJsonbDB() DB {
	db, err := sql.Open("sqlite", "./db.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	return jsonbDB{
		db: db,
	}
}

func (db jsonbDB) CreateTable() error {
	_, err := db.db.Exec(`CREATE TABLE test(x BLOB)`)
	return err
}

func (db jsonbDB) DropTable() error {
	_, err := db.db.Exec(`DROP TABLE IF EXISTS test`)
	return err
}

func (db jsonbDB) Insert(doc string) error {
	_, err := db.db.Exec(`INSERT INTO test values(jsonb(?))`, doc)
	return err
}
func (db jsonbDB) Query() (string, error) {
	row := db.db.QueryRow(`SELECT unmarshal(x) FROM test`)

	var res string
	if err := row.Scan(&res); err != nil {
		return "", err
	}

	return res, nil
}
