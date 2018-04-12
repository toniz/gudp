# Mysql DB Sharding

## GUDP DB Configure File:  
```json
{
    "db_t_account_0" : 
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
    "db_t_account_1" : 
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
    "ACCOUNT_t_user_sharding":
    {
        "sql" : "SELECT user_id, user_name, type FROM t_user;",
	"sharding": {"dbseq": ""},
        "db" : "db_t_account_$dbseq$"
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
```

## BUILD
go build client.go 
Then client will got the data result:
```
Result Like: db_t_account_0(accountdb): SELECT user_id, user_name, type FROM t_user
Result Like: db_t_account_1(image):SELECT user_id, user_name, type FROM t_user
error: code = Unknown desc = Error: [db_t_account_2]Not Found In Configure..

```




