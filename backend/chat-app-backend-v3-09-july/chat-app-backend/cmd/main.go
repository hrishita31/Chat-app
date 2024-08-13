package main

import (
	"fmt"

	"github.com/lpernett/godotenv"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("error occurred", zap.Error(err))
		return
	}

	fx.New(Modules).Run()
}
