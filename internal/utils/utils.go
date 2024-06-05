package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"flag"
	"log"

	"github.com/caarlos0/env"
)

type Flags struct {
	FlagAddr      string
	FlagDBAddr    string
	FlagSecretKey string
	SecretKey     [16]byte
}

type ServerENV struct {
	Address   string `env:"RUN_ADDRESS"`
	DBAddress string `env:"DATABASE_URI"`
	SecretKey string `env:"SECRET_KEY"`
}

func ShaData(result string, key string) string {
	b := []byte(result)
	shakey := []byte(key)
	// создаём новый hash.Hash, вычисляющий контрольную сумму SHA-256
	h := hmac.New(sha256.New, shakey)
	// передаём байты для хеширования
	h.Write(b)
	// вычисляем хеш
	hash := h.Sum(nil)
	sha := base64.URLEncoding.EncodeToString(hash)
	return string(sha)
}

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

func IsCardNumberValid(number uint64) bool {
	return (number%10+cardChecksum(number/10))%10 == 0
}

func cardChecksum(number uint64) uint64 {
	var luhn uint64
	for i := 0; number > 0; i++ {
		cur := number % 10

		if i%2 == 0 { // even
			cur = cur * 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}

		luhn += cur
		number = number / 10
	}
	return luhn % 10
}
