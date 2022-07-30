package sqlite

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"social-network/pkg/logger"

	_ "github.com/mattn/go-sqlite3"
	migrate "github.com/rubenv/sql-migrate"
)

type DB struct {
	db *sql.DB
}

func (dataBase DB) Connect(path string) (*sql.DB, error) {
	//DB initialization and connection
	storage, err := sql.Open("sqlite3", "file:"+path+"?_foreign_keys=on")
	if err != nil {
		logger.ErrorLogger.Println(err)
		return nil, err
	}
	if err := storage.Ping(); err != nil {
		return nil, err
	}
	dataBase.db = storage
	logger.InfoLogger.Println("Connect to DB successfully")

	if err := dataBase.createStorage(); err != nil {
		return nil, err
	}
	logger.InfoLogger.Println("DataBase created successfully")

	migrations := &migrate.FileMigrationSource{
		Dir: "pkg/sqlite/migrations",
	}

	n, err := migrate.Exec(storage, "sqlite3", migrations, migrate.Up)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return nil, err
	}
	fmt.Printf("Applied %d migrations!\n", n)

	return storage, nil
}

func (dataBase DB) createStorage() error {
	createStorage, err := ioutil.ReadFile("./createTables.sql")
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	_, err = dataBase.db.Exec(string(createStorage))
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	return nil
}
