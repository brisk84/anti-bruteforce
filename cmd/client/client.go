package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	pb "anti-bruteforce/internal/server/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

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

func main() {
	var (
		login string
		pass  string
		ip    string
	)

	flag.Parse()
	if len(flag.Args()) < 1 {
		fmt.Println("Usage: ab-client l login pass ip")
		return
	}

	conn, err := grpc.Dial("localhost:4242", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewAntiBruteforceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if flag.Arg(0) == "l" {
		if len(flag.Args()) < 4 {
			fmt.Println("Usage: ab-client l login pass ip")
			return
		}
		login = flag.Arg(1)
		pass = flag.Arg(2)
		ip = flag.Arg(3)
		ret, err := Login(ctx, c, login, pass, ip)
		CheckRetErr(ret, err)
	}

	if flag.Arg(0) == "r" {
		if len(flag.Args()) < 4 {
			fmt.Println("Usage: ab-client r login pass ip")
			return
		}
		login = flag.Arg(1)
		pass = flag.Arg(2)
		ip = flag.Arg(3)
		ret, err := Reset(ctx, c, login, pass, ip)
		CheckRetErr(ret, err)
	}

	if flag.Arg(0) == "ab" {
		if len(flag.Args()) < 2 {
			fmt.Println("Usage: ab-client ab ip")
			return
		}
		ip = flag.Arg(1)
		ret, err := AddToBlackList(ctx, c, ip)
		CheckRetErr(ret, err)
	}

	if flag.Arg(0) == "db" {
		if len(flag.Args()) < 2 {
			fmt.Println("Usage: ab-client db ip")
			return
		}
		ip = flag.Arg(1)
		ret, err := DelFromBlackList(ctx, c, ip)
		CheckRetErr(ret, err)
	}

	if flag.Arg(0) == "aw" {
		if len(flag.Args()) < 2 {
			fmt.Println("Usage: ab-client aw ip")
			return
		}
		ip = flag.Arg(1)
		ret, err := AddToWhiteList(ctx, c, ip)
		CheckRetErr(ret, err)
	}

	if flag.Arg(0) == "dw" {
		if len(flag.Args()) < 2 {
			fmt.Println("Usage: ab-client dw ip")
			return
		}
		ip = flag.Arg(1)
		ret, err := DelFromWhiteList(ctx, c, ip)
		CheckRetErr(ret, err)
	}
}
