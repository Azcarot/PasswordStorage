package storage

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/Azcarot/PasswordStorage/internal/cfg"

	_ "modernc.org/sqlite"
)

// LiteConn - интерфейсный тип для работы с пользователем на клиенте
type LiteConn interface {
	CreateTablesForGoKeeper()
	CreateNewUser(ctx context.Context, data RegisterRequest) error
	GetSecretKey(data string) error
}

// ClientUserStorage - тип для хранения пользовтеля на клиенте
type ClientUserStorage struct {
	Storage PgxStorage
	DB      *sql.Conn
	Data    any
}

// SQLLiteStore - тип для пользователького хранилища sqlite
type SQLLiteStore struct {
	DB *sql.DB
}

type liteConnTime struct {
	attempts          int
	timeBeforeAttempt int
}

// SyncReq - тип для передачи хешей в запросах на синхронизацию со стороны клиента
type SyncReq struct {
	BankCard string `json:"bank"`
	TextData string `json:"text"`
	FileData string `json:"file"`
	LoginPw  string `json:"lpw"`
}

// LiteDB - подключение к SQLite
var LiteDB *sql.DB

// LiteST - интерфейс для работы с пользователем на клиенте
var LiteST LiteConn

// ServURL - url сервера
var ServURL string

// AuthToken - токен авторизации
var AuthToken string

// Secret - секрет на клиенте
var Secret [16]byte

// SyncClientHashes - хэши данных пользователя для запросов на синхру с сервером
var SyncClientHashes SyncReq

// MakeLiteConn - реализация подключения к SQLlite
func MakeLiteConn(db *sql.DB) LiteConn {
	return &SQLLiteStore{
		DB: db,
	}
}

// NewLiteConn - подключение к SQLlite
func NewLiteConn(f cfg.Flags) error {
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
	ServURL = "http://" + f.FlagAddr
	return nil
}

func connectToLiteDB(f cfg.Flags) error {
	var err error
	ps := fmt.Sprintf(f.FlagDBAddr)
	LiteDB, err = sql.Open("sqlite", ps)
	LiteST = MakeLiteConn(LiteDB)
	BCLiteS = NewBCLiteStorage(BCLiteS, LiteDB)
	TLiteS = NewTLiteStorage(TLiteS, LiteDB)
	LPWLiteS = NewLPLiteStorage(LPWLiteS, LiteDB)
	FLiteS = NewFLiteStorage(FLiteS, LiteDB)
	return err
}

// CreateTablesForGoKeeper - создание таблиц для всех типов данных пользователя на клиенте
func (store SQLLiteStore) CreateTablesForGoKeeper() {
	ctx := context.Background()

	mut.Lock()
	defer mut.Unlock()
	queryForFun := `DROP TABLE IF EXISTS users CASCADE`
	store.DB.ExecContext(ctx, queryForFun)
	query := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY, 
		login TEXT NOT NULL, 
		password TEXT NOT NULL,
		secret TEXT NOT NULL,
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
		file_path TEXT NOT NULL,
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

// GenerateSecretKey - генерация секретного ключа для пользователя
func GenerateSecretKey(data RegisterRequest) string {
	// Combine login and password
	combined := data.Login + ":" + data.Password

	// Use a secret key for HMAC, this should be a constant or derived securely
	hmacKey := []byte(data.Login)

	// Create HMAC using SHA-256
	h := hmac.New(sha256.New, hmacKey)
	h.Write([]byte(combined))

	// Get the HMAC hash as a hexadecimal string
	secretKey := hex.EncodeToString(h.Sum(nil))

	return secretKey
}

// GetSecretKey - получение секрета текущего пользователя
func (store SQLLiteStore) GetSecretKey(login string) error {
	query := `SELECT DISTINCT secret
	FROM users
	WHERE login = $1`

	rows := store.DB.QueryRow(query, login)
	var skey string
	if err := rows.Scan(&skey); err != nil {
		return err
	}
	byteSlice := []byte(skey)

	// Copy the relevant portion of the byte slice into the result array
	copy(Secret[:], byteSlice)

	return nil
}

// CreateNewUser - создание нового пользователя на клиенте
func (store SQLLiteStore) CreateNewUser(ctx context.Context, data RegisterRequest) error {

	mut.Lock()
	defer mut.Unlock()
	tx, err := store.DB.BeginTx(ctx, nil)
	if err != nil {

		return err
	}
	newKey := GenerateSecretKey(data)
	date := time.Now().Format(time.RFC3339)
	_, err = store.DB.ExecContext(ctx, `INSERT into users (login, password, secret, created) 
	values ($1, $2, $3, $4);`,
		data.Login, data.Password, newKey, date)

	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}
	return err
}
