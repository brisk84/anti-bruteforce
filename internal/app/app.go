package app

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/time/rate"
)

type Bucket struct {
	login string
	pass  string
	ip    string
	lim   *rate.Limiter
}

type App struct {
	lim        *rate.Limiter
	buckets    map[string]Bucket
	whiteList  map[string]struct{}
	blackList  map[string]struct{}
	loginLimit int
	passLimit  int
	ipLimit    int
}

func New(loginLimit, passLimit, ipLimit int) *App {
	return &App{
		lim:        rate.NewLimiter(rate.Every(10*time.Second), 1),
		buckets:    make(map[string]Bucket),
		whiteList:  make(map[string]struct{}),
		blackList:  make(map[string]struct{}),
		loginLimit: loginLimit,
		passLimit:  passLimit,
		ipLimit:    ipLimit,
	}
}

func (a *App) Login(ctx context.Context, li LoginInfo) error {
	fmt.Println("In app login:", li.Login, li.Password, li.Ip)

	allow := false
	if _, ok := a.whiteList[li.Ip]; ok {
		return nil
	}
	if _, ok := a.blackList[li.Ip]; ok {
		return fmt.Errorf("ok=false")
	}

	if _, ok := a.buckets[li.Login]; !ok {
		bucket := Bucket{login: li.Login, lim: rate.NewLimiter(rate.Every(1*time.Minute), a.loginLimit)}
		a.buckets[li.Login] = bucket
	}
	if _, ok := a.buckets[li.Password]; !ok {
		bucket := Bucket{pass: li.Password, lim: rate.NewLimiter(rate.Every(1*time.Minute), a.passLimit)}
		a.buckets[li.Password] = bucket
	}
	if _, ok := a.buckets[li.Ip]; !ok {
		bucket := Bucket{ip: li.Ip, lim: rate.NewLimiter(rate.Every(1*time.Minute), a.ipLimit)}
		a.buckets[li.Ip] = bucket
	}
	allow = a.buckets[li.Login].lim.Allow() && a.buckets[li.Password].lim.Allow() && a.buckets[li.Ip].lim.Allow()
	if !allow {
		return fmt.Errorf("ok=false")
	}

	return nil
}

func (a *App) Reset(ctx context.Context, li LoginInfo) error {
	fmt.Println("In app reset:", li.Login, li.Password, li.Ip)
	delete(a.buckets, li.Login)
	delete(a.buckets, li.Password)
	delete(a.buckets, li.Ip)
	return nil
}

func (a *App) AddToBlackList(ctx context.Context, ni NetworkInfo) error {
	fmt.Println("In app AddToBlackList:", ni.Ip)
	a.blackList[ni.Ip] = struct{}{}
	return nil
}

func (a *App) DelFromBlackList(ctx context.Context, ni NetworkInfo) error {
	fmt.Println("In app DelFromBlackList:", ni.Ip)
	delete(a.blackList, ni.Ip)
	return nil
}

func (a *App) AddToWhiteList(ctx context.Context, ni NetworkInfo) error {
	fmt.Println("In app AddToWhiteList:", ni.Ip)
	a.whiteList[ni.Ip] = struct{}{}
	return nil
}

func (a *App) DelFromWhiteList(ctx context.Context, ni NetworkInfo) error {
	fmt.Println("In app DelFromWhiteList:", ni.Ip)
	delete(a.whiteList, ni.Ip)
	return nil
}
