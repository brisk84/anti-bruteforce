package server

import (
	"context"
	"fmt"
	"net"

	app "anti-bruteforce/internal/app"
	pb "anti-bruteforce/internal/server/api"
	"google.golang.org/grpc"
)

type Application interface {
	Login(context.Context, app.LoginInfo) error
	Reset(context.Context, app.LoginInfo) error
	AddToBlackList(context.Context, app.NetworkInfo) error
	DelFromBlackList(context.Context, app.NetworkInfo) error
	AddToWhiteList(context.Context, app.NetworkInfo) error
	DelFromWhiteList(context.Context, app.NetworkInfo) error
}

type Server struct {
	pb.UnimplementedAntiBruteforceServer
	addr       string
	grpcServer *grpc.Server
	app        Application
}

func NewServer(app Application, addr string) *Server {
	return &Server{app: app, addr: addr}
}

func (s *Server) Start(ctx context.Context) error {
	s.grpcServer = grpc.NewServer()

	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	pb.RegisterAntiBruteforceServer(s.grpcServer, s)
	err = s.grpcServer.Serve(lis)
	<-ctx.Done()
	return err
}

func (s *Server) Stop(ctx context.Context) error {
	s.grpcServer.Stop()
	return nil
}

func (s *Server) Login(ctx context.Context, li *pb.LoginInfo) (*pb.Error, error) {
	var err error
	fmt.Println("Login:", li.Login, li.Password, li.Ip)
	err1 := s.app.Login(ctx, app.LoginInfo{Login: li.Login, Password: li.Password, IP: li.Ip})
	if err1 == nil {
		return &pb.Error{Code: 200, Info: "ok=true"}, err
	}
	return &pb.Error{Code: 404, Info: "ok=false"}, err
}

func (s *Server) Reset(ctx context.Context, li *pb.LoginInfo) (*pb.Error, error) {
	var err error
	fmt.Println("Reset:", li.Login, li.Password, li.Ip)
	err1 := s.app.Reset(ctx, app.LoginInfo{Login: li.Login, Password: li.Password, IP: li.Ip})
	_ = err1
	return &pb.Error{Code: 404, Info: "Not found"}, err
}

func (s *Server) AddToBlackList(ctx context.Context, ni *pb.NetworkInfo) (*pb.Error, error) {
	s.app.AddToBlackList(ctx, app.NetworkInfo{IP: ni.Ip})
	return &pb.Error{Code: 200, Info: "Ok"}, nil
}

func (s *Server) DelFromBlackList(ctx context.Context, ni *pb.NetworkInfo) (*pb.Error, error) {
	s.app.DelFromBlackList(ctx, app.NetworkInfo{IP: ni.Ip})
	return &pb.Error{Code: 200, Info: "Ok"}, nil
}

func (s *Server) AddToWhiteList(ctx context.Context, ni *pb.NetworkInfo) (*pb.Error, error) {
	s.app.AddToWhiteList(ctx, app.NetworkInfo{IP: ni.Ip})
	return &pb.Error{Code: 200, Info: "Ok"}, nil
}

func (s *Server) DelFromWhiteList(ctx context.Context, ni *pb.NetworkInfo) (*pb.Error, error) {
	s.app.DelFromWhiteList(ctx, app.NetworkInfo{IP: ni.Ip})
	return &pb.Error{Code: 200, Info: "Ok"}, nil
}
