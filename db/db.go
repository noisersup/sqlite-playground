package db

type DB interface {
	CreateTable() error
	DropTable() error
	Insert(doc string) error
	Query() (string, error)
}
