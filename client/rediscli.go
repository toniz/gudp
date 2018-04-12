package main

import (
	"flag"
	"fmt"
	"log"
	"reflect"

	pb "bbwhat.net/gudserver/interface"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "", "The file containning the CA root cert file")
	serverAddr         = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	serverHostOverride = flag.String("server_host_override", "www.bbwhat.net", "The server name use to verify the hostname returned by TLS handshake")
)

func main() {
	flag.Parse()
	var opts []grpc.DialOption
	if *tls {
		creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
		if err != nil {
			log.Fatalf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	fmt.Println(conn.GetState())
	defer conn.Close()
	client := pb.NewUnifiedDataClient(conn)

	fmt.Println(reflect.TypeOf(conn))
	{
		req := pb.Query{
			Engine: "redis",
			Ident:  "srv_1",
			Opt:    "set",
			Params: map[string]string{
				"key":   "uid_199",
				"value": "10",
			},
		}
		stream, err := client.DBCommit(context.Background(), &req)
		if err != nil {
			log.Fatalf("%v.ListFeatures(_) = _, %v", client, err)
		}
		fmt.Println(stream)
	}
	{
		req := pb.Query{
			Engine: "redis",
			Ident:  "srv_2",
			Opt:    "get",
			Params: map[string]string{
				"key": "uid_199",
			},
		}
		stream, err := client.DBCommit(context.Background(), &req)
		if err != nil {
			log.Fatalf("%v.ListFeatures(_) = _, %v", client, err)
		}
		fmt.Println(stream)
	}

}
