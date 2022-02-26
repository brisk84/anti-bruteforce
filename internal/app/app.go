package app

import (
	"anti-bruteforce/internal/logger"
	"context"
	"fmt"
	"net"
	"time"

	"golang.org/x/time/rate"
)

type Bucket struct {
	login   string
	pass    string
	ip      string
	lim     *rate.Limiter
	expired time.Time
}

type App struct {
	lim        *rate.Limiter
	buckets    map[string]Bucket
	whiteList  map[string]struct{}
	blackList  map[string]struct{}
	loginLimit int
	passLimit  int
	ipLimit    int
	ticker     *time.Ticker
	logg       *logger.Logger
	testMode   bool
}

func New(ctx context.Context, logger *logger.Logger, loginLimit, passLimit, ipLimit int, testMode bool) *App {
	bucketTimeout := 1 * time.Minute
	app := App{
		lim:        rate.NewLimiter(rate.Every(bucketTimeout), 1),
		buckets:    make(map[string]Bucket),
		whiteList:  make(map[string]struct{}),
		blackList:  make(map[string]struct{}),
		loginLimit: loginLimit,
		passLimit:  passLimit,
		ipLimit:    ipLimit,
		ticker:     time.NewTicker(bucketTimeout),
		logg:       logger,
		testMode:   testMode,
	}
	go app.CleanBuckets(ctx)
	return &app
}

func (a *App) CleanBuckets(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-a.ticker.C:
			for k, v := range a.buckets {
				if time.Until(v.expired) < 1*time.Nanosecond {
					delete(a.buckets, k)
				}
			}
		}
	}
}

func (a *App) Login(ctx context.Context, li LoginInfo) error {
	a.logg.Info("In app login: " + li.Login + " " + li.Password + " " + li.IP)

	for k := range a.whiteList {
		_, mask, _ := net.ParseCIDR(k)
		if mask.Contains(net.ParseIP(li.IP)) {
			return nil
		}
	}
	for k := range a.blackList {
		_, mask, _ := net.ParseCIDR(k)
		if mask.Contains(net.ParseIP(li.IP)) {
			return fmt.Errorf("ok=false")
		}
	}

	timeInt := 1 * time.Minute
	if a.testMode {
		timeInt = 10 * time.Second
	}
	timeExp := time.Now().Add(timeInt)
	if _, ok := a.buckets[li.Login]; !ok {
		bucket := Bucket{login: li.Login, lim: rate.NewLimiter(rate.Every(timeInt), a.loginLimit), expired: timeExp}
		a.buckets[li.Login] = bucket
	}
	if _, ok := a.buckets[li.Password]; !ok {
		bucket := Bucket{pass: li.Password, lim: rate.NewLimiter(rate.Every(timeInt), a.passLimit), expired: timeExp}
		a.buckets[li.Password] = bucket
	}
	if _, ok := a.buckets[li.IP]; !ok {
		bucket := Bucket{ip: li.IP, lim: rate.NewLimiter(rate.Every(timeInt), a.ipLimit), expired: timeExp}
		a.buckets[li.IP] = bucket
	}
	allow := a.buckets[li.Login].lim.Allow() && a.buckets[li.Password].lim.Allow() && a.buckets[li.IP].lim.Allow()
	if !allow {
		return fmt.Errorf("ok=false")
	}

	return nil
}

func (a *App) Reset(ctx context.Context, li LoginInfo) error {
	a.logg.Info("In app reset: " + li.Login + " " + li.Password + " " + li.IP)
	delete(a.buckets, li.Login)
	delete(a.buckets, li.Password)
	delete(a.buckets, li.IP)
	return nil
}

func (a *App) AddToBlackList(ctx context.Context, ni NetworkInfo) error {
	a.logg.Info("In app AddToBlackList: " + ni.IP)
	a.blackList[ni.IP] = struct{}{}
	return nil
}

func (a *App) DelFromBlackList(ctx context.Context, ni NetworkInfo) error {
	a.logg.Info("In app DelFromBlackList: " + ni.IP)
	delete(a.blackList, ni.IP)
	return nil
}

func (a *App) AddToWhiteList(ctx context.Context, ni NetworkInfo) error {
	a.logg.Info("In app AddToWhiteList: " + ni.IP)
	a.whiteList[ni.IP] = struct{}{}
	return nil
}

func (a *App) DelFromWhiteList(ctx context.Context, ni NetworkInfo) error {
	a.logg.Info("In app DelFromWhiteList: " + ni.IP)
	delete(a.whiteList, ni.IP)
	return nil
}
