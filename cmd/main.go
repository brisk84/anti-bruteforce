package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	app "anti-bruteforce/internal/app"
	server "anti-bruteforce/internal/server"
)

func main() {
	ab := app.New(10, 100, 1000)
	grpcServer := server.NewServer(ab, "localhost:4242")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := grpcServer.Stop(ctx); err != nil {
			fmt.Println(err)
		}
	}()
	fmt.Println("Server started")

	if err := grpcServer.Start(ctx); err != nil {
		fmt.Println(err)
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
