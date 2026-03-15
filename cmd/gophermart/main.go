package main

import (
	"github.com/klyakssa/go-diplom.git/internal/app"
	"github.com/klyakssa/go-diplom.git/internal/config"
)

func main() {
	config := config.InitConfiguration()

	app.Run(config)
}
