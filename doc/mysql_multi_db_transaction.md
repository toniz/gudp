# Mysql Multi DB Transaction

## GUDP DB Configure File:  
```json
{
    "db_t_account_w" : 
    {
        "DBName" : "accountdb",
        "DBUser" : "account_rw",
        "DBPass" : "123456",
        "ConnString" : "127.0.0.1:4000",
        "ConnMaxIdleTime" : 60,
        "ConnTimeout" : 5,
        "ConnMaxCnt" : 100,
        "ConnMaxLifetime" : 3600,
        "ConnEncoding" : "utf8,utf8mb4"
    },
    "db_t_image_w" : 
    {
        "DBName" : "image",
        "DBUser" : "uds",
        "DBPass" : "uds1234",
        "ConnString" : "127.0.0.1:4000",
        "ConnMaxIdleTime" : 60,
        "ConnTimeout" : 5,
        "ConnMaxCnt" : 100,
        "ConnMaxLifetime" : 3600,
        "ConnEncoding" : "utf8,utf8mb4"
    }
}

```

## GUDP SQL Configure File:  
```json
{
    "ACCOUNT_t_user_insert_transaction" : 
    {
	"sqlgroup": 
	    [
            {
                "sql" : "INSERT INTO t_user(user_id, user_name, type) VALUES($id$, $name$ ,$type$);",
                "noquote": {"id":""},
                "check":   {"id": "^\\d+$"},
                "db" : "db_t_account_w"
            },
            { 
                "sql" : "INSERT INTO t_images(id, name, image) VALUES($id$, $name$ ,$image$);",
                "noquote": {"id":""},
                "check":   {"id": "^\\d+$"},
                "db" : "db_t_image_w"
            }
	    ]
    }
}

```

## Then The GRPC Clien.go Call Like this:  

```go
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
}
```

## BUILD
go build client.go 
Then client will got the data like follow sql :
```
Result Like: db_t_account_0(accountdb): INSERT INTO t_user(user_id, user_name, type) VALUES(11000, "uet" ,"23");
Result Like: db_t_image_w(image): INSERT INTO t_images(id, name, image) VALUES(10002, "uae" ,"111");
```



