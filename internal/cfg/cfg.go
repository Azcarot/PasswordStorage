// Package cfg - обработка флагов и шифрование данных пользователя
package cfg

import (
	"flag"
	"log"

	"github.com/caarlos0/env"
)

// Flags - тип для хранения флагов/переменных окружения
type Flags struct {
	FlagAddr      string
	FlagDBAddr    string
	FlagSecretKey string
	SecretKey     [16]byte
}

// ServerENV -тип переменных окружения
type ServerENV struct {
	Address   string `env:"RUN_ADDRESS"`
	DBAddress string `env:"DATABASE_URI"`
	SecretKey string `env:"SECRET_KEY"`
}

// ParseFlagsAndENV - получение значений флагов и переменных окружения
func ParseFlagsAndENV() Flags {
	var Flag Flags
	flag.StringVar(&Flag.FlagAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&Flag.FlagDBAddr, "d", "", "address for db")
	flag.StringVar(&Flag.FlagSecretKey, "s", "", "secret")
	flag.Parse()
	var envcfg ServerENV
	err := env.Parse(&envcfg)
	if err != nil {
		log.Fatal(err)
	}

	if len(envcfg.Address) > 0 {
		Flag.FlagAddr = envcfg.Address
	}
	if len(envcfg.DBAddress) > 0 {
		Flag.FlagDBAddr = envcfg.DBAddress
	}

	if len(envcfg.SecretKey) > 0 {
		Flag.FlagSecretKey = envcfg.SecretKey
	}

	//Ключ делаем 16-байтным
	var byteArray [16]byte

	byteSlice := []byte(Flag.FlagSecretKey)

	if len(byteSlice) > 16 {

		copy(byteArray[:], byteSlice[:16])
	} else {

		copy(byteArray[:], byteSlice)
	}
	Flag.SecretKey = byteArray

	return Flag
}
