package db

import (
	"database/sql"
	"log"
)

type jsonDB struct {
	db *sql.DB
}

func NewJsonDB() jsonDB {
	db, err := sql.Open("sqlite", "./db.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	return jsonDB{
		db: db,
	}
}

func (db jsonDB) CreateTable() error {
	_, err := db.db.Exec(`CREATE TABLE test(x TEXT)`)
	return err
}

func (db jsonDB) DropTable() error {
	_, err := db.db.Exec(`DROP TABLE IF EXISTS test`)
	return err
}

func (db jsonDB) Insert(doc string) error {
	_, err := db.db.Exec(`INSERT INTO test VALUES(?)`, doc)
	return err
}

func (db jsonDB) Query() (string, error) {
	row := db.db.QueryRow(`SELECT x FROM test`)

	var res string
	if err := row.Scan(&res); err != nil {
		return "", err
	}

	return res, nil
}
