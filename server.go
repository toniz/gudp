package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"strings"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	//_ "google.golang.org/grpc/grpclog/glogger"

	pb "bbwhat.net/gudserver/interface"
	"bbwhat.net/gudserver/mysql"
	"bbwhat.net/gudserver/redis"
)

var (
	tls      = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile = flag.String("cert_file", "", "The TLS cert file")
	keyFile  = flag.String("key_file", "", "The TLS key file")
	port     = flag.Int("port", 10000, "The server port")
)

type UnifiedData struct {
}

type UnifiedDataer interface {
	DoCommit(ctx context.Context, req *pb.Query) (*pb.Response, error)
}

var mysqler UnifiedDataer
var rediser UnifiedDataer

// DBCommit returns the query result .
func (s *UnifiedData) DBCommit(ctx context.Context, req *pb.Query) (res *pb.Response, err error) {
	grpclog.Infoln("req: ", req)

	switch strings.ToLower(req.Engine) {
	case "mysql":
		res, err = mysqler.DoCommit(ctx, req)
	case "redis":
		res, err = rediser.DoCommit(ctx, req)
	default:
		msg := "Not Include [" + req.Engine + "] Engine"
		grpclog.Warningf(msg)
		err = errors.New(msg)
	}
	grpclog.Infoln("res: ", res)
	return res, err
}

func newServer() (*UnifiedData, error) {
	s := &UnifiedData{}
	var err error
	if mysqler, err = mysql.NewMysqlDataService(); err != nil {
		grpclog.Errorf("Mysql Data Service Initialization Failed: %v", err)
	}

	if rediser, err = redis.NewRedisDataService(); err != nil {
		grpclog.Errorf("Redis Data Service Initialization Failed: %v", err)
	}
	return s, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		grpclog.Errorf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	if *tls {
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			grpclog.Errorf("Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	grpcServer := grpc.NewServer(opts...)
	if s, err := newServer(); err == nil {
		pb.RegisterUnifiedDataServer(grpcServer, s)
		grpcServer.Serve(lis)
	} else {
		grpclog.Errorf("GrpcDataServer Start Failed. %v", err)
	}
}
