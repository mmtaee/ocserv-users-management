package bootstrap

import (
	"api/pkg/config"
	"api/pkg/database"
	"api/pkg/routing"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Serve(debug bool, host string, port int) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic recovered: %v", r)
		}
	}()

	cfg := config.NewConfig(debug, host, port)

	database.Connect(cfg)
	Migrate()

	defer database.Close()

	go routing.Serve(cfg)

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
