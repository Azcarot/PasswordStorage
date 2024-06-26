package main

import (
	"fmt"
	"os"

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
			panic(err)
		}

	}

	storage.LiteConn.CreateTablesForGoKeeper(storage.LiteST)
	p := face.MakeTeaProg()
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
