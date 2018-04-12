package config

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"google.golang.org/grpc/grpclog"
)

var (
	defaultPath = os.Getenv("GOPATH") + "/src/github.com/toniz/gudp/"
)

func NewMysqlConfig() (*MysqlConfig, error) {
	s := &MysqlConfig{DC: make(MysqlDBs), SC: make(MysqlSQLs)}
	err := s.LoadConfigure()
	return s, err
}

func NewRedisConfig() (*RedisConfig, error) {
	s := &RedisConfig{RC: make(RedisSRVs)}
	err := s.LoadConfigure()
	return s, err
}

type CfgLoader interface{}

// Read and parse all json files under path
func loadFromFile(path string, l CfgLoader) error {
	if files, err := ioutil.ReadDir(path); err != nil {
		return err
	} else {
		for _, file := range files {
			jsonFile, err := os.Open(path + file.Name())
			if err != nil {
				return err
			}
			defer jsonFile.Close()

			jsonStream, _ := ioutil.ReadAll(jsonFile)
			err = json.Unmarshal(jsonStream, l)
			if err != nil {
				grpclog.Errorf("SQL Configure File[%s] loading Failed. %v", file.Name(), err)
				return err
			}
		}
	}
	//	grpclog.Infof("Loading Configure From Json File: %v", l)

	return nil
}

func loadFromDB() error {
	grpclog.Infof("Loading Configure From DB")
	return nil
}

func loadFromETCD() error {
	grpclog.Infof("Loading Configure ETCD")
	return nil
}
