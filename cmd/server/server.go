package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	signal "os/signal"
	"syscall"
	"time"

	app "anti-bruteforce/internal/app"
	server "anti-bruteforce/internal/server"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	loginLimit := flag.Int("loginLimit", 10, "Login limit")
	passLimit := flag.Int("passLimit", 100, "Pass limit")
	ipLimit := flag.Int("ipLimit", 1000, "IP limit")
	flag.Parse()

	ab := app.New(ctx, *loginLimit, *passLimit, *ipLimit)
	grpcServer := server.NewServer(ab, "localhost:4242")

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
