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

	// Create a context for cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rcl := NewRedisClient()
	_, err := rcl.redisCL.Ping(ctx).Result()
	if err != nil {
		fmt.Println("err: pingingngn----", err.Error())
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
		fmt.Printf("Starting server at %s", port)

		if err := srv.ListenAndServe(); err != nil {
			fmt.Println(err)
			log.Fatal(err)
		}
	}()

	// Wait for termination signal
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	<-signalCh

	// Cancel the context and exit
	cancel()
}
