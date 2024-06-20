package storage

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Azcarot/PasswordStorage/internal/cfg"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type CtxKey string

var mut sync.Mutex

// UserLoginCtxKey - ключ для хранения логина в конткесте запросов
const UserLoginCtxKey CtxKey = "userLogin"

// EncryptionCtxKey - ключ для хранения уникального sercret пользователя в конткесте запросов
const EncryptionCtxKey CtxKey = "secret"

// UserData - тип для хранения пользовательских данных
// для авторизации/регистрации
type UserData struct {
	Login    string
	Password string
	Date     string
}

// LoginData - тип данных пользователя формата логин/пароль
type LoginData struct {
	ID       int    `json:"id"`
	Login    string `json:"login"`
	Password string `json:"pwd"`
	Comment  string `json:"comment"`
	User     string
	Date     string
	Str      string
}

// TextData - тип данных пользователя формата текстовые данные
type TextData struct {
	ID      int    `json:"id"`
	Text    string `json:"text"`
	Comment string `json:"comment"`
	User    string
	Date    string
	Str     string
}

// FileData - тип данных пользователя формата "файл"
type FileData struct {
	ID       int    `json:"id"`
	FileName string `json:"name"`
	Path     string `json:"path"`
	Data     string `json:"data"`
	Comment  string `json:"comment"`
	User     string
	Date     string
	Str      string
}

// CardData - тип данных пользователя формата банковских карт
type BankCardData struct {
	ID         int    `json:"id"`
	CardNumber string `json:"number"`
	Cvc        string `json:"cvc"`
	ExpDate    string `json:"exp"`
	FullName   string `json:"fullname"`
	Comment    string `json:"comment"`
	User       string
	Date       string
	Str        string
}

// GetReqData - формат данных для запросов Get и Search
type GetReqData struct {
	ID  string `json:"id"`
	Str string `json:"str"`
}

// PgxConn интерфейс подключения/работы с пользователем в бд
type PgxConn interface {
	CreateTablesForGoKeeper()
	CreateNewUser(ctx context.Context, data UserData) error
	CheckUserExists(ctx context.Context, data UserData) (bool, error)
	CheckUserPassword(ctx context.Context, data UserData) (bool, error)
}

// PgxStorage - интерфейс работы с хранилищами данных пользователя
type PgxStorage interface {
	CreateNewRecord(ctx context.Context) error
	UpdateRecord(ctx context.Context) error
	GetRecord(ctx context.Context) (any, error)
	GetAllRecords(ctx context.Context) (any, error)
	DeleteRecord(ctx context.Context) error
	HashDatabaseData(ctx context.Context) (string, error)
	AddData(data any) error
	GetData() any
}

// SQLStore - тип для подключения к постгресс
type SQLStore struct {
	DB *pgx.Conn
}

// Storage - тип для работы с пользователем в БД
type Storage struct {
	Storage PgxStorage
	DB      *pgx.Conn
	Data    any
}

// BankCardStorage - тип для работы с банковскими картами в postgress
type BankCardStorage struct {
	Storage PgxStorage
	DB      *pgx.Conn
	Data    BankCardData
}

// LoginPWStorage - тип для работы с логин/пароль в postgress
type LoginPwStorage struct {
	Storage PgxStorage
	DB      *pgx.Conn
	Data    LoginData
}

// TextStorage - тип для работы с текстовыми данными в postgress
type TextStorage struct {
	Storage PgxStorage
	DB      *pgx.Conn
	Data    TextData
}

// FileStorage - тип для работы с файлами в postgress
type FileStorage struct {
	Storage PgxStorage
	DB      *pgx.Conn
	Data    FileData
}

// SyncData - тип для хранения хешей данных, для запросов синхронизации
type SyncData struct {
	BankCardHash string
	TextDataHash string
	FileData     string
	LoginData    string
}

// SyncServerHashes - хэши данных пользователя с сервера
var SyncServerHashes SyncData

// UserLoginPw - данные пользователя (логин/пароль)
var UserLoginPw UserData

// DB подключение к бд для работы с пользователем
var DB *pgx.Conn

// ST - интерфейс работы с пользователем
var ST PgxConn

// BCST - интерфейс работы с банковскими картами в pgx
var BCST PgxStorage

// LPST - интерфейс работы с данными типа логин/пароль в pgx
var LPST PgxStorage

// TST - интерфейс работы с текстовыми данными в pgx
var TST PgxStorage

// FST - интерфейс работы с файловыми данными в pgx
var FST PgxStorage

// ErrNoLogin - ошибка при отсутствии логина пользователя в контексте
var ErrNoLogin = fmt.Errorf("no login in context")

// ErrNoKey - ошибка при отсутсвии секрета в контексте
var ErrNoKey = fmt.Errorf("no key in context")

type pgxConnTime struct {
	attempts          int
	timeBeforeAttempt int
}

// NewPwStorage - реализация хранилища для пользователя
func NewPwStorage(storage PgxStorage, db *pgx.Conn) *Storage {
	return &Storage{
		Storage: storage,
		DB:      db,
	}
}

// NewBCStorage - реализация хранилища банковских карт
func NewBCStorage(storage PgxStorage, db *pgx.Conn) *BankCardStorage {
	return &BankCardStorage{
		Storage: storage,
		DB:      db,
	}
}

// NewLPSTtorage - реализация хранилища логина/пароля
func NewLPSTtorage(storage PgxStorage, db *pgx.Conn) *LoginPwStorage {
	return &LoginPwStorage{
		Storage: storage,
		DB:      db,
	}
}

// NewTSTtorage - реализация хранилища текстовых данных
func NewTSTtorage(storage PgxStorage, db *pgx.Conn) *TextStorage {
	return &TextStorage{
		Storage: storage,
		DB:      db,
	}
}

// NewFSTtorage - реализация хранилища файлов
func NewFSTtorage(storage PgxStorage, db *pgx.Conn) *FileStorage {
	return &FileStorage{
		Storage: storage,
		DB:      db,
	}
}

// MakeConn - реализация хранилища пользователей
func MakeConn(db *pgx.Conn) PgxConn {
	return &SQLStore{
		DB: db,
	}
}

// NewConn - подключение к БД, с последющим созданием и реализацией всех хранилищ для сервера
func NewConn(f cfg.Flags) error {
	var err error
	var attempts pgxConnTime
	attempts.attempts = 3
	attempts.timeBeforeAttempt = 1
	err = connectToDB(f)
	for err != nil {
		//если ошибка связи с бд, то это не эскпортируемый тип, отличный от PgError
		var pqErr *pgconn.PgError
		if errors.Is(err, pqErr) {
			return err

		}
		if attempts.attempts == 0 {
			return err
		}
		times := time.Duration(attempts.timeBeforeAttempt)
		time.Sleep(times * time.Second)
		attempts.attempts--
		attempts.timeBeforeAttempt += 2
		err = connectToDB(f)

	}
	return nil
}

func connectToDB(f cfg.Flags) error {
	var err error
	ps := fmt.Sprintf(f.FlagDBAddr)
	DB, err = pgx.Connect(context.Background(), ps)
	ST = MakeConn(DB)
	BCST = NewBCStorage(BCST, DB)
	LPST = NewLPSTtorage(LPST, DB)
	TST = NewTSTtorage(TST, DB)
	FST = NewFSTtorage(FST, DB)
	return err
}

// AddData - добавление данных в хранилище
func (store *BankCardStorage) AddData(data any) error {
	newdata, ok := data.(BankCardData)
	if !ok {
		return fmt.Errorf("error while asserting data to bank card type")
	}
	store.Data = newdata
	return nil
}

// GetData - получение данных из хранилища
func (store *BankCardStorage) GetData() any {

	newdata := store.Data

	return newdata
}

func (store *BankCardLiteStorage) AddData(data any) error {

	newdata, ok := data.(BankCardData)
	if !ok {
		return fmt.Errorf("error while asserting data to bank card type")
	}
	store.Data = newdata
	return nil
}

func (store *BankCardLiteStorage) GetData() any {

	newdata := store.Data

	return newdata
}

func (store *LPWLiteStorage) AddData(data any) error {

	newdata, ok := data.(LoginData)
	if !ok {
		return fmt.Errorf("error while asserting data to bank card type")
	}
	store.Data = newdata
	return nil
}

func (store *LPWLiteStorage) GetData() any {

	newdata := store.Data

	return newdata
}

func (store *FileLiteStorage) AddData(data any) error {

	newdata, ok := data.(FileData)
	if !ok {
		return fmt.Errorf("error while asserting data to file type")
	}
	store.Data = newdata
	return nil
}

func (store *FileLiteStorage) GetData() any {

	newdata := store.Data

	return newdata
}

func (store *TextLiteStorage) AddData(data any) error {

	newdata, ok := data.(TextData)
	if !ok {
		return fmt.Errorf("error while asserting data to text type")
	}
	store.Data = newdata
	return nil
}

func (store *TextLiteStorage) GetData() any {

	newdata := store.Data

	return newdata
}

func (store *LoginPwStorage) AddData(data any) error {

	newdata, ok := data.(LoginData)
	if !ok {
		return fmt.Errorf("error while asserting data to login/pw type")
	}
	store.Data = newdata
	return nil
}

func (store *LoginPwStorage) GetData() any {

	newdata := store.Data

	return newdata
}

func (store *TextStorage) AddData(data any) error {

	newdata, ok := data.(TextData)
	if !ok {
		return fmt.Errorf("error while asserting data to text type")
	}
	store.Data = newdata
	return nil
}

func (store *TextStorage) GetData() any {

	newdata := store.Data

	return newdata
}

func (store *FileStorage) AddData(data any) error {

	newdata, ok := data.(FileData)
	if !ok {
		return fmt.Errorf("error while asserting data to file type")
	}
	store.Data = newdata
	return nil
}

func (store *FileStorage) GetData() any {

	newdata := store.Data

	return newdata
}

// CheckDBConnection - проверка соединения с БД
func CheckDBConnection() http.Handler {
	checkConnection := func(res http.ResponseWriter, req *http.Request) {

		err := DB.Ping(context.Background())
		result := (err == nil)
		if result {
			res.WriteHeader(http.StatusOK)
		} else {
			res.WriteHeader(http.StatusInternalServerError)
		}

	}
	return http.HandlerFunc(checkConnection)
}

// CreateTablesForGoKeeper - создание таблиц для всех типов данных в бд сервера
func (store SQLStore) CreateTablesForGoKeeper() {
	ctx := context.Background()
	mut.Lock()
	defer mut.Unlock()
	queryForFun := `DROP TABLE IF EXISTS users CASCADE`
	store.DB.Exec(ctx, queryForFun)
	query := `CREATE TABLE IF NOT EXISTS users (
		id SERIAL NOT NULL PRIMARY KEY, 
		login TEXT NOT NULL, 
		password TEXT NOT NULL,
		created text )`

	_, err := store.DB.Exec(ctx, query)

	if err != nil {

		log.Printf("Error %s when creating user table", err)

	}
	queryForFun = `DROP TABLE IF EXISTS bank_card CASCADE`
	store.DB.Exec(ctx, queryForFun)
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
	_, err = store.DB.Exec(ctx, query)

	if err != nil {

		log.Printf("Error %s when creating bank_card table", err)

	}

	queryForFun = `DROP TABLE IF EXISTS login_pw CASCADE`
	store.DB.Exec(ctx, queryForFun)
	query = `CREATE TABLE IF NOT EXISTS login_pw(
		id SERIAL NOT NULL PRIMARY KEY,
		login TEXT NOT NULL,
		pw TEXT NOT NULL,
		comment TEXT,
		username TEXT NOT NULL,
		created TEXT
	)`
	_, err = store.DB.Exec(ctx, query)

	if err != nil {

		log.Printf("Error %s when creating login_pw table", err)

	}

	queryForFun = `DROP TABLE IF EXISTS text_data CASCADE`
	store.DB.Exec(ctx, queryForFun)
	query = `CREATE TABLE IF NOT EXISTS text_data(
		id SERIAL NOT NULL PRIMARY KEY,
		text TEXT NOT NULL,
		comment TEXT,
		username TEXT NOT NULL,
		created TEXT
	)`
	_, err = store.DB.Exec(ctx, query)

	if err != nil {

		log.Printf("Error %s when creating login_pw table", err)

	}

	queryForFun = `DROP TABLE IF EXISTS file_data CASCADE`
	store.DB.Exec(ctx, queryForFun)
	query = `CREATE TABLE IF NOT EXISTS file_data(
		id SERIAL NOT NULL PRIMARY KEY,
		file_name TEXT NOT NULL,
		file_path TEXT NOT NULL,
		data TEXT NOT NULL,
		comment TEXT,
		username TEXT NOT NULL,
		created TEXT
	)`
	_, err = store.DB.Exec(ctx, query)

	if err != nil {

		log.Printf("Error %s when creating file_data table", err)

	}

}
