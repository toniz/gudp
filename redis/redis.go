package redis

import (
	"context"
	"errors"
	"strings"

	"github.com/go-redis/redis"
	"google.golang.org/grpc/grpclog"

	pb "bbwhat.net/toniz/gudp/interface"
	conf "github.com/toniz/gudp/config"
)

type RedisHandle map[string]*redis.Client
type RedisDataService struct {
	rdh RedisHandle
}

// NewMysqlDataService creates a MysqlDataService with the provided GrpcDataServer.
// Call loadConfigure To Load Configure Data
func NewRedisDataService() (*RedisDataService, error) {
	s := &RedisDataService{rdh: make(RedisHandle)}
	l, err := conf.NewRedisConfig()
	if err != nil || len(l.RC) == 0 {
		return nil, err
	}

	for k, v := range l.RC {
		// Connect To Redis
		client := redis.NewClient(&redis.Options{
			Addr:     v.ConnString,
			Password: v.DBPass, // no password set
			DB:       v.DBSeq,  // use default DB
		})

		_, err := client.Ping().Result()
		if err != nil {
			grpclog.Errorf("Redis[%s] Connect Failed [%v]: %v", k, v, err)
		}
		s.rdh[k] = client
	}
	return s, nil
}

// Exec Redis Command.
func (s *RedisDataService) DoCommit(ctx context.Context, req *pb.Query) (*pb.Response, error) {

	var res pb.Response
	var err error

	if rdh := s.rdh[req.Ident]; rdh != nil {
		switch strings.ToUpper(req.Opt) {
		case "SET":
			err = rdh.Set(req.Params["key"], req.Params["value"], 0).Err()
			if err == nil {
				res.Count = 1
			}
		case "GET":
			var val string
			key := req.Params["key"]
			val, err = rdh.Get(key).Result()
			if err == redis.Nil {
				msg := "key [" + key + "]does not exist"
				err = errors.New(msg)
				grpclog.Infof(msg)
			}

			if err == nil {
				// Now do something with the data.
				// Here we just print each column as a string.
				data := make(map[string]string)
				data[key] = val
				var row pb.RowData
				row.Fields = data
				res.Rows = append(res.Rows, &row)
			}
		default:
			msg := "Redis Opt [" + req.Opt + "] not definded."
			err = errors.New(msg)
			grpclog.Warningf(msg)
		}

		if err != nil {
			grpclog.Warningf("Get Redis Data Failed: %v", err)
		}
	} else {
		msg := "Error: [" + req.Ident + "]No Such Redis Server"
		err = errors.New(msg)
		grpclog.Errorf(msg)
	}
	return &res, err
}
