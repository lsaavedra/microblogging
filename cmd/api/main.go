package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"microblogging/db"
	"microblogging/internal/config"
	"microblogging/internal/handler"
)

func main() {
	go beforeShutdown()

	envConf := config.Load(os.Getenv("ENV"))

	dbInstance := db.Init(envConf)

	router := handler.RouterWithHandlers(dbInstance)

	initServer(router, envConf)
}

func initServer(router *gin.Engine, envConf *config.Config) {
	err := router.Run(fmt.Sprintf(":%s", envConf.ListeningPort))
	if err != nil {
		panic(err)
	}
}

func beforeShutdown() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println("for some reason, I'm going down. Trying to release and close all the resources. Bye!: ", sig)
		os.Exit(1)
	}()
}
