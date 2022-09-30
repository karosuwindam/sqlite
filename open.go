package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type sqliteConfig struct {
	filepass string
	db       *sql.DB
}

func Setup(filepath string) sqliteConfig {
	return sqliteConfig{filepass: filepath}
}

func (t *sqliteConfig) Open() error {
	var err error
	t.db, err = sql.Open("sqlite3", t.filepass)
	return err
}

func (t *sqliteConfig) Close() error {
	return t.db.Close()

}

func (t *sqliteConfig) ReturnFilePass() string { return t.filepass }
