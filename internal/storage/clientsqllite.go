package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Azcarot/PasswordStorage/internal/utils"
	_ "modernc.org/sqlite"
)

type LiteConn interface {
	CreateTablesForGoKeeper()
}

type ClientUserStorage struct {
	Storage PgxStorage
	DB      *sql.Conn
	Data    any
}

type SQLLiteStore struct {
	DB *sql.DB
}

type liteConnTime struct {
	attempts          int
	timeBeforeAttempt int
}

var LiteDB *sql.DB
var LiteST LiteConn
var ServURL string
var AuthToken string

func MakeLiteConn(db *sql.DB) LiteConn {
	return &SQLLiteStore{
		DB: db,
	}
}

func NewLiteConn(f utils.Flags) error {
	var err error
	var attempts liteConnTime
	attempts.attempts = 3
	attempts.timeBeforeAttempt = 1
	err = connectToLiteDB(f)
	for err != nil {

		if attempts.attempts == 0 {
			return err
		}
		times := time.Duration(attempts.timeBeforeAttempt)
		time.Sleep(times * time.Second)
		attempts.attempts -= 1
		attempts.timeBeforeAttempt += 2
		err = connectToLiteDB(f)

	}
	ServURL = f.FlagAddr
	return nil
}

func connectToLiteDB(f utils.Flags) error {
	var err error
	ps := fmt.Sprintf(f.FlagDBAddr)
	LiteDB, err = sql.Open("sqlite", ps)
	LiteST = MakeLiteConn(LiteDB)
	return err
}

func (store SQLLiteStore) CreateTablesForGoKeeper() {
	ctx := context.Background()

	mut.Lock()
	defer mut.Unlock()
	queryForFun := `DROP TABLE IF EXISTS users CASCADE`
	store.DB.ExecContext(ctx, queryForFun)
	query := `CREATE TABLE IF NOT EXISTS users (
		id SERIAL NOT NULL PRIMARY KEY, 
		login TEXT NOT NULL, 
		password TEXT NOT NULL,
		created text )`

	_, err := store.DB.ExecContext(ctx, query)

	if err != nil {

		log.Printf("Error %s when creating user table", err)

	}
	queryForFun = `DROP TABLE IF EXISTS bank_card CASCADE`
	store.DB.ExecContext(ctx, queryForFun)
	query = `CREATE TABLE IF NOT EXISTS bank_card(
		id SERIAL NOT NULL PRIMARY KEY,
		card_number TEXT NOT NULL,
		cvc TEXT NOT NULL,
		exp_date TEXT,
		full_name TEXT NOT NULL,
		comment TEXT,
		username TEXT NOT NULL,
		created TEXT
	)`
	_, err = store.DB.ExecContext(ctx, query)

	if err != nil {

		log.Printf("Error %s when creating bank_card table", err)

	}

	queryForFun = `DROP TABLE IF EXISTS login_pw CASCADE`
	store.DB.ExecContext(ctx, queryForFun)
	query = `CREATE TABLE IF NOT EXISTS login_pw(
		id SERIAL NOT NULL PRIMARY KEY,
		login TEXT NOT NULL,
		pw TEXT NOT NULL,
		comment TEXT,
		username TEXT NOT NULL,
		created TEXT
	)`
	_, err = store.DB.ExecContext(ctx, query)

	if err != nil {

		log.Printf("Error %s when creating login_pw table", err)

	}

	queryForFun = `DROP TABLE IF EXISTS text_data CASCADE`
	store.DB.ExecContext(ctx, queryForFun)
	query = `CREATE TABLE IF NOT EXISTS text_data(
		id SERIAL NOT NULL PRIMARY KEY,
		text TEXT NOT NULL,
		comment TEXT,
		username TEXT NOT NULL,
		created TEXT
	)`
	_, err = store.DB.ExecContext(ctx, query)

	if err != nil {

		log.Printf("Error %s when creating login_pw table", err)

	}

	queryForFun = `DROP TABLE IF EXISTS file_data CASCADE`
	store.DB.ExecContext(ctx, queryForFun)
	query = `CREATE TABLE IF NOT EXISTS file_data(
		id SERIAL NOT NULL PRIMARY KEY,
		file_name TEXT NOT NULL,
		data TEXT NOT NULL,
		comment TEXT,
		username TEXT NOT NULL,
		created TEXT
	)`
	_, err = store.DB.ExecContext(ctx, query)

	if err != nil {

		log.Printf("Error %s when creating file_data table", err)

	}

}
