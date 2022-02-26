/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	pb "anti-bruteforce/internal/server/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "ab-client <command> <params>",
	Short: "CLI interface for anti-bruteforce service",
	Long:  `CLI interface for anti-bruteforce service`,
}

var (
	client pb.AntiBruteforceClient
	ctx    context.Context
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	viper.SetConfigName("client_config")
	viper.AddConfigPath("configs")
	cfg := NewConfig()
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(&cfg)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	conn, err := grpc.Dial(net.JoinHostPort(cfg.Server.Host, cfg.Server.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client = pb.NewAntiBruteforceClient(conn)

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err = rootCmd.Execute()
	if err != nil {
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func Login(ctx context.Context, c pb.AntiBruteforceClient, login, pass, ip string) (string, error) {
	li := pb.LoginInfo{Login: login, Password: pass, Ip: ip}
	resp, err := c.Login(ctx, &li)
	if err != nil {
		return "", err
	}
	return resp.Info, nil
}

func Reset(ctx context.Context, c pb.AntiBruteforceClient, login, pass, ip string) (string, error) {
	li := pb.LoginInfo{Login: login, Password: pass, Ip: ip}
	resp, err := c.Reset(ctx, &li)
	if err != nil {
		return "", err
	}
	return resp.Info, nil
}

func AddToBlackList(ctx context.Context, c pb.AntiBruteforceClient, ip string) (string, error) {
	ni := pb.NetworkInfo{Ip: ip}
	resp, err := c.AddToBlackList(ctx, &ni)
	if err != nil {
		return "", err
	}
	return resp.Info, nil
}

func DelFromBlackList(ctx context.Context, c pb.AntiBruteforceClient, ip string) (string, error) {
	ni := pb.NetworkInfo{Ip: ip}
	resp, err := c.DelFromBlackList(ctx, &ni)
	if err != nil {
		return "", err
	}
	return resp.Info, nil
}

func AddToWhiteList(ctx context.Context, c pb.AntiBruteforceClient, ip string) (string, error) {
	ni := pb.NetworkInfo{Ip: ip}
	resp, err := c.AddToWhiteList(ctx, &ni)
	if err != nil {
		return "", err
	}
	return resp.Info, nil
}

func DelFromWhiteList(ctx context.Context, c pb.AntiBruteforceClient, ip string) (string, error) {
	ni := pb.NetworkInfo{Ip: ip}
	resp, err := c.DelFromWhiteList(ctx, &ni)
	if err != nil {
		return "", err
	}
	return resp.Info, nil
}

func CheckRetErr(ret string, err error) {
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ret)
}
