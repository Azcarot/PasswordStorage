package storage

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Azcarot/PasswordStorage/internal/utils"
	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const SecretKey string = "super-secret"

type MyCustomClaims struct {
	jwt.MapClaims
}

type CtxKey string

var mut sync.Mutex

const UserLoginCtxKey CtxKey = "userLogin"
const DBCtxKey CtxKey = "dbConn"
const EncryptionCtxKey CtxKey = "secret"

type UserData struct {
	Login    string
	Password string
	Date     string
}

type LoginData struct {
	ID       string `json:"id"`
	Login    string `json:"login"`
	Password string `json:"pwd"`
	Comment  string `json:"comment"`
	User     string
	Date     string
	Str      string
}

type LoginResponse struct {
	Login    string `json:"login"`
	Password string `json:"pwd"`
	Comment  string `json:"comment"`
}

type TextData struct {
	ID      string `json:"id"`
	Text    string `json:"text"`
	Comment string `json:"comment"`
	User    string
	Date    string
	Str     string
}

type TextResponse struct {
	Text    string `json:"text"`
	Comment string `json:"comment"`
}

type FileData struct {
	ID       string `json:"id"`
	FileName string `json:"name"`
	Data     string `json:"data"`
	Comment  string `json:"comment"`
	User     string
	Date     string
	Str      string
}

type FileResponse struct {
	FileName string `json:"name"`
	Data     string `json:"data"`
	Comment  string `json:"comment"`
}

type BankCardData struct {
	ID         string `json:"id"`
	CardNumber string `json:"number"`
	Cvc        string `json:"cvc"`
	ExpDate    string `json:"exp"`
	FullName   string `json:"fullname"`
	Comment    string `json:"comment"`
	User       string
	Date       string
	Str        string
}

type BankCardResponse struct {
	CardNumber string `json:"number"`
	Cvc        string `json:"cvc"`
	ExpDate    string `json:"exp"`
	FullName   string `json:"fullname"`
	Comment    string `json:"comment"`
}

type GetReqData struct {
	ID  string `json:"id"`
	Str string `json:"str"`
}

type PgxConn interface {
	CreateTablesForGoKeeper()
	CreateNewUser(ctx context.Context, data UserData) error
	CheckUserExists(data UserData) (bool, error)
	CheckUserPassword(ctx context.Context, data UserData) (bool, error)
}

type PgxStorage interface {
	CreateNewRecord(ctx context.Context) error
	UpdateRecord(ctx context.Context) error
	GetRecord(ctx context.Context) (any, error)
	GetAllRecords(ctx context.Context) (any, error)
	SearchRecord(ctx context.Context) (any, error)
	DeleteRecord(ctx context.Context) error
	AddData(data any) error
}

type SQLStore struct {
	DB *pgx.Conn
}

type Storage struct {
	Storage PgxStorage
	DB      *pgx.Conn
	Data    any
}

type BankCardStorage struct {
	Storage PgxStorage
	DB      *pgx.Conn
	Data    BankCardData
}

type LoginPwStorage struct {
	Storage PgxStorage
	DB      *pgx.Conn
	Data    LoginData
}

type TextStorage struct {
	Storage PgxStorage
	DB      *pgx.Conn
	Data    TextData
}

type FileStorage struct {
	Storage PgxStorage
	DB      *pgx.Conn
	Data    FileData
}

var DB *pgx.Conn
var ST PgxConn
var PST PgxStorage
var BCST PgxStorage
var LPST PgxStorage
var TST PgxStorage
var FST PgxStorage
var errTimeout = fmt.Errorf("timeout exceeded")
var ErrNoLogin = fmt.Errorf("no login in context")
var ErrNoKey = fmt.Errorf("no key in context")

type pgxConnTime struct {
	attempts          int
	timeBeforeAttempt int
}

type WithdrawRequest struct {
	OrderNumber string  `json:"order"`
	Amount      float64 `json:"sum"`
}

type WithdrawResponse struct {
	OrderNumber string  `json:"order"`
	Amount      float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}

func NewPwStorage(storage PgxStorage, db *pgx.Conn) *Storage {
	return &Storage{
		Storage: storage,
		DB:      db,
	}
}

func NewBCStorage(storage PgxStorage, db *pgx.Conn) *BankCardStorage {
	return &BankCardStorage{
		Storage: storage,
		DB:      db,
	}
}

func NewLPSTtorage(storage PgxStorage, db *pgx.Conn) *LoginPwStorage {
	return &LoginPwStorage{
		Storage: storage,
		DB:      db,
	}
}

func NewTSTtorage(storage PgxStorage, db *pgx.Conn) *TextStorage {
	return &TextStorage{
		Storage: storage,
		DB:      db,
	}
}

func NewFSTtorage(storage PgxStorage, db *pgx.Conn) *FileStorage {
	return &FileStorage{
		Storage: storage,
		DB:      db,
	}
}

func MakeConn(db *pgx.Conn) PgxConn {
	return &SQLStore{
		DB: db,
	}
}

func NewConn(f utils.Flags) error {
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
		attempts.attempts -= 1
		attempts.timeBeforeAttempt += 2
		err = connectToDB(f)

	}
	return nil
}

func connectToDB(f utils.Flags) error {
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

func (store *BankCardStorage) AddData(data any) error {

	newdata, ok := data.(BankCardData)
	if !ok {
		return fmt.Errorf("error while asserting data to bank card type")
	}
	store.Data = newdata
	return nil
}

func (store *LoginPwStorage) AddData(data any) error {

	newdata, ok := data.(LoginData)
	if !ok {
		return fmt.Errorf("error while asserting data to login/pw type")
	}
	store.Data = newdata
	return nil
}

func (store *TextStorage) AddData(data any) error {

	newdata, ok := data.(TextData)
	if !ok {
		return fmt.Errorf("error while asserting data to text type")
	}
	store.Data = newdata
	return nil
}

func (store *FileStorage) AddData(data any) error {

	newdata, ok := data.(FileData)
	if !ok {
		return fmt.Errorf("error while asserting data to file type")
	}
	store.Data = newdata
	return nil
}

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

func CypherData(ctx context.Context, data string) (string, error) {
	keyData, ok := ctx.Value(EncryptionCtxKey).([16]byte)
	if !ok {
		return "", ErrNoKey
	}
	key := keyData[:]
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(data))

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// decrypt decrypts ciphertext using the given key and returns the plaintext
func Dechypher(ctx context.Context, data string) (string, error) {
	keyData, ok := ctx.Value(EncryptionCtxKey).([16]byte)
	if !ok {
		return "", ErrNoKey
	}
	key := keyData[:]
	ciphertextBytes, err := base64.URLEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(ciphertextBytes) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext слишком короткий")
	}

	iv := ciphertextBytes[:aes.BlockSize]
	ciphertextBytes = ciphertextBytes[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertextBytes, ciphertextBytes)

	return string(ciphertextBytes), nil
}
