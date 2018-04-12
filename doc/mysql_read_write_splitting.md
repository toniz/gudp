# Mysql Read Write Splitting

## DB Grant 
```
GRANT ALL PRIVILEGES ON accountdb.* TO 'account_rw'@'127.0.0.1' IDENTIFIED BY '123456' WITH GRANT OPTION;
GRANT SELECT  ON accountdb.* TO 'account_r'@'127.0.0.1' IDENTIFIED BY '123456' WITH GRANT OPTION;
```

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
    "db_t_account_r" : 
    {
        "DBName" : "accountdb",
        "DBUser" : "account_r",
        "DBPass" : "123456",
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
    "ACCOUNT_t_user_insert" :
    {
        "sql" : "INSERT INTO t_user(user_id, user_name, type) VALUES($id$, $name$ ,$type$);",
        "noquote": {"id":""},
        "check":   {"id": "^\\d+$"},
        "db" : "db_t_account_w"
    },
    "ACCOUNT_t_user_update" :
    {
        "sql" : "UPDATE t_user SET user_name = $name$ WHERE USER_ID = $id$;",
        "db" : "db_t_gpsbox_w"
    },
    "ACCOUNT_t_user_delete" :
    {
        "sql" : "DELETE FROM t_user WHERE user_id = $id$",
        "db" : "db_t_gpsbox_w"
    },
    "ACCOUNT_t_user_insert_multi" :
    {
        "sql" : "REPLACE INTO $table_name$(user_id, user_name, type) VALUES $values$",
        "noquote": {"table_name":"", "values":""},
        "noescape":{"values":""},
        "check":   {"values": "^.*$"},
        "db" : "db_t_gpsbox_w"
    },
    "ACCOUNT_t_user_select_by_uid" :
    {
        "sql" : "SELECT user_id, user_name, type FROM t_user WHERE user_id>=$limit_start$ ORDER BY user_id ASC LIMIT $limit_end$ ;",
        "noquote" : {"limit_start":"", "limit_end":""},
        "db" : "db_t_account_r"
    },
    "ACCOUNT_t_user_select_by_uids" :
    {
        "sql" : "SELECT user_id, user_name, type FROM t_user WHERE user_id in ($condition$)",
        "noquote": {"condition":""},
        "noescape":{"condition":""},
        "check": {"condition": "^.*$"},
        "db" : "db_t_gpsbox_w"
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
}
```

## BUILD
go build client.go 
Then client will got the data like follow sql :
```
SELECT user_id, user_name, type FROM t_user WHERE user_id>=100 ORDER BY user_id ASC LIMIT 10 ;
```


