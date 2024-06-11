// Package main - главный серверный модуль. Запуск осуществляется с флагами -
// -d - авторизация к бд, -a -адрес сервера, -s - секрет
package main

import (
	"context"
	"log"
	"net/http"

	"github.com/Azcarot/PasswordStorage/internal/router"
	"github.com/Azcarot/PasswordStorage/internal/storage"
	"github.com/Azcarot/PasswordStorage/internal/utils"
)

func main() {
	flag := utils.ParseFlagsAndENV()
	if flag.FlagDBAddr == "" {
		log.Fatal("Missing required flag -d : DataBase address\n")
	}
	err := storage.NewConn(flag)
	if err != nil {
		panic(err)
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
