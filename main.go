package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
)

func main() {
	var identifier string = os.Getenv("IDENTIFIER")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rcl := NewRedisClient()

	for {
		_, err := rcl.Ping(ctx)
		if err == nil {
			break
		}

		time.Sleep(3 * time.Second)
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

	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err)

		return
	}

	clienset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)

		return
	}

	lock, err := resourcelock.New(
		resourcelock.LeasesResourceLock,
		"default",
		"lease-lock",
		clienset.CoreV1(),
		clienset.CoordinationV1(),
		resourcelock.ResourceLockConfig{
			Identity: identifier,
		},
	)
	if err != nil {
		log.Fatal(err)

		return
	}

	lconf := leaderelection.LeaderElectionConfig{
		Lock:          lock,
		LeaseDuration: 10 * time.Second,
		RenewDeadline: 5 * time.Second,
		RetryPeriod:   2 * time.Second,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(ctx context.Context) {
				go processData(ctx, identifier)
			},
			OnStoppedLeading: func() {
				fmt.Println(identifier, " has stopped leading")
			},
			OnNewLeader: func(identity string) {
				if identity == identifier {
					fmt.Println(identity, " is still the leader")

					return
				}

				fmt.Println(identity, " is now the new leader")
			},
		},
	}

	lElector, err := leaderelection.NewLeaderElector(lconf)
	if err != nil {
		log.Fatal(err)

		return
	}

	lElector.Run(ctx)

	// Wait for termination signal
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	<-signalCh

	if err := srv.Shutdown(ctx); err != nil {
		rcl.redisCL.Shutdown(ctx)
	}
}
