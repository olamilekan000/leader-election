package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	var identifier string = os.Getenv("IDENTIFIER")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rcl := NewRedisClient()

	_, err := rcl.Ping(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Redis is ready. Starting the application.")

	router := gin.Default()
	port := ":8881"

	AddRoutes(router)

	srv := &http.Server{
		Addr:    port,
		Handler: router,
	}

	go func() {
		fmt.Println("Starting server at " + port)

		if err := srv.ListenAndServe(); err != nil {
			log.Fatalln(err)
		}
	}()

	go processData(ctx, identifier)

	// Wait for termination signal
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	<-signalCh

	if err := srv.Shutdown(ctx); err != nil {
		rcl.redisCL.Shutdown(ctx)
	}
}
