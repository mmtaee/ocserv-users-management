package bootstrap

import (
	"context"
	"fmt"
	"log"
	"ocserv-bakend/pkg/config"
	"ocserv-bakend/pkg/database"
	"ocserv-bakend/pkg/routing"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Serve(debug bool) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic recovered: %v", r)
		}
	}()

	config.Init(debug)
	database.Connect(debug)
	Migrate()

	defer database.Close()

	go routing.Serve(debug)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Second signal received, forcing exit")
		os.Exit(1)
	}()

	sig := <-quit
	fmt.Println()
	log.Printf("signal %v received", sig)
	log.Println("initiating shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	routing.Shutdown(ctx)
}
