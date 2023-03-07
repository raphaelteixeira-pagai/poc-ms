package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golangsugar/chatty"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/raphaelteixeira-pagai/poc-ms/internal/domain/handlers"
	"github.com/raphaelteixeira-pagai/poc-ms/internal/domain/repository"
	"github.com/raphaelteixeira-pagai/poc-ms/internal/domain/services"
	"github.com/raphaelteixeira-pagai/poc-ms/pkg/database"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if err := godotenv.Load(); err != nil {
		chatty.Fatal(err.Error())
	}

	conn, _ := database.NewPostgres(1)
	repo := repository.NewWalletRepository(conn)
	srv := services.NewWalletService(repo)

	app := gin.Default()
	router := app.Group("/")
	handlers.RegisterWalletRoutes(srv, router)

	// setup server
	server := &http.Server{
		Addr:              fmt.Sprintf(":%v", os.Getenv("PORT")),
		Handler:           app,
		WriteTimeout:      time.Second * 15,
		ReadTimeout:       time.Second * 15,
		ReadHeaderTimeout: time.Second * 15,
		IdleTimeout:       time.Second * 60,
	}

	// start server and wait for OS signal
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			chatty.FatalErr(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	_ = server.Shutdown(ctx)
	chatty.Info("shutting down! 👋")
}
