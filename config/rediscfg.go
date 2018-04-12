package config

import (
	"errors"
	"flag"
)

var (
	redisSrvPath  = flag.String("redis_srv_path", defaultPath+"conf/redis/srv/", "Redis SRV configure path")
	redisLoadType = flag.Int("redis_load_type", 1, "Redis Load Type: 1 [Json File]; 2 [DB]; 3 [ETCD]")
)

type RedisSRV struct {
	ConnString string
	DBPass     string
	DBSeq      int
}

type RedisSRVs map[string]RedisSRV

type RedisConfig struct {
	RC RedisSRVs
}

// loadConfigure loads configure from a JSON file Directory.
func (s *RedisConfig) LoadConfigure() error {
	switch *redisLoadType {
	case 1:

		if err := loadFromFile(*redisSrvPath, &s.RC); err != nil {
			return err
		}

	case 2:
		if err := loadFromDB(); err != nil {
			return err
		}
	case 3:
		if err := loadFromETCD(); err != nil {
			return err
		}
	default:
		return errors.New("Error: Redis Config Load Type")
	}
	return nil
}
