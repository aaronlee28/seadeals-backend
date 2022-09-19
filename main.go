package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"seadeals-backend/config"
	"seadeals-backend/db"
	"seadeals-backend/redis"
	"seadeals-backend/server"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error cannot load file local.env")
	}
	config.Reset()

	err = db.Connect()
	if err != nil {
		fmt.Println("failed to connect to DB")
		return
	}

	redis.Setup()
	server.Init()
}
