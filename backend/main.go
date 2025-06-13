package main

import (
	"fmt"
	"log"
	_ "ocserv-bakend/docs"
	"ocserv-bakend/pkg/commands"
	"ocserv-bakend/pkg/config"
	"ocserv-bakend/pkg/database"
	"ocserv-bakend/pkg/routing"
	"os"
	"os/signal"
	"syscall"
)

// @title Ocserv User management Example Api
// @version 1.0
// @description This is a sample Ocserv User management Api server.
// @BasePath /api
func main() {

	debug := os.Getenv("DEBUG") == "true"

	config.Init(debug)
	database.Connect(debug)
	commands.Migrate()
	defer database.Close()

	go func() {
		routing.Serve(debug)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	sig := <-quit

	fmt.Println()
	log.Printf("signal %v received\n", sig)
	log.Println("initiating shutdown...")

	routing.Shutdown()

}
