package storage

import (
	"context"
	"errors"
	"fmt"
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
const OrderNumberCtxKey CtxKey = "orderNumber"
const DBCtxKey CtxKey = "dbConn"

type UserData struct {
	Login    string
	Password string
	Date     string
}

type LoginData struct {
	Login    string `json:"login"`
	Password string `json:"pwd"`
	Comment  string `json:"comment"`
	User     string
}

type BankCardData struct {
	CardNumber uint64 `json:"number"`
	Cvc        int    `json:"cvc"`
	ExpDate    string `json:"exp"`
	FullName   string `json:"fullname"`
	Comment    string `json:"comment"`
	User       string
}

type BankCardResponse struct {
	CardNumber uint64 `json:"number"`
	Cvc        int    `json:"cvc"`
	ExpDate    string `json:"exp"`
	FullName   string `json:"fullname"`
	Comment    string `json:"comment"`
}

type BalanceResponce struct {
	Accrual   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type PgxConn interface {
	CreateTablesForGoKeeper()
	CreateNewUser(ctx context.Context, data UserData) error
	CheckUserExists(data UserData) (bool, error)
	CheckUserPassword(ctx context.Context, data UserData) (bool, error)
}

type PgxStorage interface {
	CreateNewRecord(ctx context.Context) error
	UpdateRecord(ctx context.Context, id int) error
	GetRecord(ctx context.Context, id int) (any, error)
	GetAllRecords(ctx context.Context) (any, error)
	SearchRecord(ctx context.Context, str string) (any, error)
	DeleteRecord(ctx context.Context, id int) error
}

type SQLStore struct {
	DB *pgx.Conn
}

type Storage struct {
	Storage PgxStorage
	DB      *pgx.Conn
	Data    any
}

var DB *pgx.Conn
var ST PgxConn
var PST PgxStorage
var errTimeout = fmt.Errorf("timeout exceeded")
var ErrNoLogin = fmt.Errorf("no login in context")

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
	PST = NewPwStorage(PST, DB)
	return err
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
		user TEXT NOT NULL, 
		comment TEXT,
		created text )`

	_, err := store.DB.Exec(ctx, query)

	if err != nil {

		log.Printf("Error %s when creating user table", err)

	}
	queryForFun = `DROP TABLE IF EXISTS bank_card CASCADE`
	store.DB.Exec(ctx, queryForFun)
	query = `CREATE TABLE IF NOT EXISTS bank_card(
		id SERIAL NOT NULL PRIMARY KEY,
		card_number BIGINT NOT NULL,
		cvc BIGINT NOT NULL,
		exp_date TEXT,
		full_name TEXT NOT NULL,
		comment TEXT,
		user TEXT NOT NULL,
		created TEXT
	)`
	_, err = store.DB.Exec(ctx, query)

	if err != nil {

		log.Printf("Error %s when creating order table", err)

	}
}
