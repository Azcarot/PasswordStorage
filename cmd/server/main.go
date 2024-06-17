// Package main - главный серверный модуль. Запуск осуществляется с флагами -
// -d - авторизация к бд, -a -адрес сервера, -s - секрет
package main

import (
	"context"
	"log"
	"net/http"

	"github.com/Azcarot/PasswordStorage/internal/cfg"
	"github.com/Azcarot/PasswordStorage/internal/router"
	"github.com/Azcarot/PasswordStorage/internal/storage"
)

func main() {
	flag := cfg.ParseFlagsAndENV()
	if flag.FlagDBAddr == "" {
		log.Fatal("Missing required flag -d : DataBase address\n")
	}
	err := storage.NewConn(flag)
	if err != nil {
		log.Fatal(err)
	}
	storage.PgxConn.CreateTablesForGoKeeper(storage.ST)
	defer storage.DB.Close(context.Background())
	r := router.MakeRouter(flag)
	server := &http.Server{
		Addr:    flag.FlagAddr,
		Handler: r,
	}
	server.ListenAndServe()

}
