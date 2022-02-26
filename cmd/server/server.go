package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	signal "os/signal"
	"syscall"
	"time"

	app "anti-bruteforce/internal/app"
	logger "anti-bruteforce/internal/logger"
	server "anti-bruteforce/internal/server"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("server_config")
	viper.AddConfigPath("configs")
	cfg := NewConfig()
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(&cfg)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	loginLimit := flag.Int("loginLimit", cfg.App.LoginLimit, "Login limit")
	passLimit := flag.Int("passLimit", cfg.App.PassLimit, "Pass limit")
	ipLimit := flag.Int("ipLimit", cfg.App.IPLimit, "IP limit")
	testMode := flag.Bool("testMode", false, "If true set timeouts to 10 second, otherwise to 1 minute")
	flag.Parse()
	fmt.Println(*testMode, *loginLimit, *passLimit, *ipLimit)

	logg := logger.New(cfg.Logger.Path, cfg.Logger.Level)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	ab := app.New(ctx, logg, *loginLimit, *passLimit, *ipLimit, *testMode)
	grpcServer := server.NewServer(ab, net.JoinHostPort(cfg.Server.Host, cfg.Server.Port), logg)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := grpcServer.Stop(ctx); err != nil {
			logg.Error(err.Error())
		}
	}()
	logg.Info("Server started")

	if err := grpcServer.Start(ctx); err != nil {
		logg.Error(err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
