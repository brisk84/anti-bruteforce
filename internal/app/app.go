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
	fmt.Println("In app login:", li.Login, li.Password, li.IP)

	var allow bool
	if _, ok := a.whiteList[li.IP]; ok {
		return nil
	}
	if _, ok := a.blackList[li.IP]; ok {
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
	if _, ok := a.buckets[li.IP]; !ok {
		bucket := Bucket{ip: li.IP, lim: rate.NewLimiter(rate.Every(1*time.Minute), a.ipLimit)}
		a.buckets[li.IP] = bucket
	}
	allow = a.buckets[li.Login].lim.Allow() && a.buckets[li.Password].lim.Allow() && a.buckets[li.IP].lim.Allow()
	if !allow {
		return fmt.Errorf("ok=false")
	}

	return nil
}

func (a *App) Reset(ctx context.Context, li LoginInfo) error {
	fmt.Println("In app reset:", li.Login, li.Password, li.IP)
	delete(a.buckets, li.Login)
	delete(a.buckets, li.Password)
	delete(a.buckets, li.IP)
	return nil
}

func (a *App) AddToBlackList(ctx context.Context, ni NetworkInfo) error {
	fmt.Println("In app AddToBlackList:", ni.IP)
	a.blackList[ni.IP] = struct{}{}
	return nil
}

func (a *App) DelFromBlackList(ctx context.Context, ni NetworkInfo) error {
	fmt.Println("In app DelFromBlackList:", ni.IP)
	delete(a.blackList, ni.IP)
	return nil
}

func (a *App) AddToWhiteList(ctx context.Context, ni NetworkInfo) error {
	fmt.Println("In app AddToWhiteList:", ni.IP)
	a.whiteList[ni.IP] = struct{}{}
	return nil
}

func (a *App) DelFromWhiteList(ctx context.Context, ni NetworkInfo) error {
	fmt.Println("In app DelFromWhiteList:", ni.IP)
	delete(a.whiteList, ni.IP)
	return nil
}
