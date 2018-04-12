package main

import (
	"flag"
	"fmt"
	"log"
	"reflect"
	"strconv"

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
			Engine: "mysql",
			Ident:  "ACCOUNT_t_user_select_by_uid",
			Params: map[string]string{
				"limit_start": "100",
				"limit_end":   "10",
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
			Engine: "mysql",
			Opt:    "2",
			Ident:  "ACCOUNT_t_user_insert_transaction",
			Group: []*pb.ParamsGroup{
				&pb.ParamsGroup{
					Params: map[string]string{
						"id":   "11000",
						"name": "uet",
						"type": "23",
					},
				},
				&pb.ParamsGroup{
					Params: map[string]string{
						"id":    "10002",
						"name":  "uae",
						"image": "111",
					},
				},
			},
		}

		stream, err := client.DBCommit(context.Background(), &req)
		if err != nil {
			log.Fatalf("%v.ListFeatures(_) = _, %v", client, err)
		}

		fmt.Println(stream)
	}

	uids := []int{100, 1003456, 2004000}
	for _, uid := range uids {
		shardnum := uid / 1000000
		req := pb.Query{
			Engine: "mysql",
			Ident:  "ACCOUNT_t_user_sharding",
			Params: map[string]string{
				"dbseq": strconv.Itoa(shardnum),
			},
		}

		stream, err := client.DBCommit(context.Background(), &req)
		if err != nil {
			log.Fatalf("%v.ListFeatures(_) = _, %v", client, err)
		}

		fmt.Println(stream)
	}
}
