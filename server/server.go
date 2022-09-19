package server

import (
	"log"
	"seadeals-backend/config"
)

func Init() {
	router := NewRouter(&RouterConfig{})
	log.Fatalln(router.Run(":" + config.Config.Port))
}
