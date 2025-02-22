package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"wallet-api/database/postgresql"
	"wallet-api/internal/app"
	"wallet-api/internal/pkg/env"
	"wallet-api/internal/pkg/logger"
	"wallet-api/internal/router"

	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	logFile := logger.Init()
	defer logFile.Close()

	db, err := postgresql.Open()
	if err != nil {
		log.Fatalf("error opening database: %v", err)
	}

	r := app.RegisterRepositories(db)
	s := app.RegisterServices(r)

	m := app.RegisterMiddlewares(s)
	h := app.RegisterHandlers(s)

	f := fiber.New()

	router.SetupRouter(f, m, h)

	go func() {
		port := env.GetWithDefault("PORT", "3000")
		log.Fatalf("error running server: %v", f.Listen(":"+port))
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig

	fmt.Println("server stopped")
	os.Exit(0)
}
