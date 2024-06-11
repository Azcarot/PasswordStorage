// Package main - главный клиентский модуль. Запуск осуществляется с флагами -
// -d - путь к файлу Sqlite, -a -адрес сервера
package main

import (
	"fmt"
	"log"

	face "github.com/Azcarot/PasswordStorage/internal/bubbleface"
	"github.com/Azcarot/PasswordStorage/internal/storage"
	"github.com/Azcarot/PasswordStorage/internal/utils"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	fmt.Printf("Build version=%s\nBuild date =%s\nBuild commit =%s\n", buildVersion, buildDate, buildCommit)
	flag := utils.ParseFlagsAndENV()
	if flag.FlagDBAddr != "" {
		err := storage.NewLiteConn(flag)
		if err != nil {
			log.Fatal(err)
		}

	}

	storage.LiteConn.CreateTablesForGoKeeper(storage.LiteST)
	p := face.MakeTeaProg()
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
