package main

import (
	"social-network/config"
	"social-network/internal/app"
	"social-network/pkg/logger"
)

func main() {
	c := config.GetConfig()

	a := new(app.App)
	err := a.Run(c.Server.Port, c.Database.DbPath)
	if err != nil {
		logger.ErrorLogger.Println(err)
		panic(err)
	}
	logger.InfoLogger.Println("Application runs")
}
