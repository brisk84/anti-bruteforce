package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/time/rate"
)

type Bucket struct {
	login string
	pass  string
	ip    string
	lim   *rate.Limiter
}

var (
	lim       *rate.Limiter
	whiteList []string
	blackList []string
	buckets   map[string]Bucket
)

func login(writer http.ResponseWriter, request *http.Request) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(writer, err.Error(), 422)
		return
	}

	var m map[string]interface{}
	err = json.Unmarshal(body, &m)
	if err != nil {
		http.Error(writer, err.Error(), 422)
		return
	}

	params := 0
	for k, v := range m {
		if (k == "login") || (k == "password") || (k == "ip") {
			if _, ok := v.(string); ok {
				params++
			}
		}
	}
	if params < 3 {
		http.Error(writer, "incorrect params", 422)
		return
	}

	allow := false
	ret := ""
	login := m["login"].(string)
	pass := m["password"].(string)
	ip := m["ip"].(string)

	if _, ok := buckets[login]; !ok {
		bucket := Bucket{login: login, lim: rate.NewLimiter(rate.Every(1*time.Minute), 2)}
		buckets[login] = bucket
	}
	if _, ok := buckets[pass]; !ok {
		bucket := Bucket{pass: pass, lim: rate.NewLimiter(rate.Every(1*time.Minute), 2)}
		buckets[pass] = bucket
	}
	if _, ok := buckets[ip]; !ok {
		bucket := Bucket{ip: ip, lim: rate.NewLimiter(rate.Every(1*time.Minute), 2)}
		buckets[ip] = bucket
	}
	allow = buckets[login].lim.Allow() && buckets[pass].lim.Allow() && buckets[ip].lim.Allow()

	ret = "ok=false"
	if allow {
		ret = "ok=true"
	}

	writer.WriteHeader(200)
	writer.Write([]byte(ret))
}

func reset(writer http.ResponseWriter, request *http.Request) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(writer, err.Error(), 422)
		return
	}

	var m map[string]interface{}
	err = json.Unmarshal(body, &m)
	if err != nil {
		http.Error(writer, err.Error(), 422)
		return
	}

	params := 0
	for k, v := range m {
		if (k == "login") || (k == "password") || (k == "ip") {
			if _, ok := v.(string); ok {
				params++
			}
		}
	}
	if params < 2 {
		http.Error(writer, "incorrect params", 422)
		return
	}

	ret := "ok"
	login := m["login"].(string)
	ip := m["ip"].(string)

	delete(buckets, login)
	delete(buckets, ip)

	writer.WriteHeader(200)
	writer.Write([]byte(ret))
}

func addBlackList(writer http.ResponseWriter, request *http.Request) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(writer, err.Error(), 422)
		return
	}

	var m map[string]interface{}
	err = json.Unmarshal(body, &m)
	// if err != nil {
	// 	http.Error(writer, err.Error(), 422)
	// 	return
	// }

	ret := ""
	if ip, ok := m["ip"]; ok {
		ret += ip.(string) + "\r"
		blackList = append(blackList, ip.(string))
	}

	host, _, err := net.SplitHostPort(request.RemoteAddr)
	if err == nil {
		ret += host
	}

	writer.WriteHeader(200)
	writer.Write([]byte(ret))
}

func delBlackList(writer http.ResponseWriter, request *http.Request) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(writer, err.Error(), 422)
		return
	}

	var m map[string]interface{}
	err = json.Unmarshal(body, &m)
	// if err != nil {
	// 	http.Error(writer, err.Error(), 422)
	// 	return
	// }

	ret := ""
	if username, ok := m["login"]; ok {
		ret += username.(string) + "\r"
	}
	if ip, ok := m["ip"]; ok {
		ret += ip.(string) + "\r"
	}

	host, _, err := net.SplitHostPort(request.RemoteAddr)
	if err == nil {
		ret += host
	}

	writer.WriteHeader(200)
	writer.Write([]byte(ret))
}

func main() {

	buckets = make(map[string]Bucket)

	lim = rate.NewLimiter(rate.Every(10*time.Second), 1)

	ctx, cancel := context.WithCancel(context.Background())

	mux := http.NewServeMux()
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/reset", reset)
	mux.HandleFunc("/add_bl", addBlackList)
	mux.HandleFunc("/del_bl", delBlackList)

	mux.HandleFunc("/h", func(writer http.ResponseWriter, request *http.Request) {
		time.Sleep(5 * time.Second)
		writer.WriteHeader(200)
	})

	httpServer := &http.Server{
		Addr:        ":80",
		Handler:     mux,
		BaseContext: func(_ net.Listener) context.Context { return ctx },
	}

	// Run server
	go func() {
		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			// it is fine to use Fatal here because it is not main gorutine
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}()

	signalChan := make(chan os.Signal, 1)

	signal.Notify(
		signalChan,
		syscall.SIGHUP,  // kill -SIGHUP XXXX
		syscall.SIGINT,  // kill -SIGINT XXXX or Ctrl+c
		syscall.SIGQUIT, // kill -SIGQUIT XXXX
	)

	<-signalChan
	log.Print("os.Interrupt - shutting down...\n")

	go func() {
		<-signalChan
		log.Fatal("os.Kill - terminating...\n")
	}()

	gracefullCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if err := httpServer.Shutdown(gracefullCtx); err != nil {
		log.Printf("shutdown error: %v\n", err)
		defer os.Exit(1)
		return
	} else {
		log.Printf("gracefully stopped\n")
	}

	cancel()

	defer os.Exit(0)
}

func main2() {
	lim = rate.NewLimiter(rate.Every(10*time.Second), 1)

	r := mux.NewRouter()
	// paths := []string{"/api1", "/api2"}
	// for _, p := range paths {
	// 	go func(path string) {
	// 		r.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
	// 			w.Write([]byte("hello world"))
	// 		})
	// 	}(p)
	// }

	// go func(path string) {
	// 	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
	// 		w.Write([]byte("hello world"))
	// 	})
	// }("/login")
	go func(path string) {
		r.HandleFunc("/login", login)
	}("/login")
	log.Fatal(http.ListenAndServe(":80", r))
}
