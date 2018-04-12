package config

import (
	"errors"
	"flag"
)

var (
	mysqlDBPath   = flag.String("mysql_db_path", defaultPath+"conf/mysql/db/", "Mysql DB configure path")
	mysqlSqlPath  = flag.String("mysql_sql_path", defaultPath+"conf/mysql/sql/", "Mysql SQL configure path")
	mysqlLoadType = flag.Int("mysql_load_type", 1, "Mysql Load Type: 1 [Json File]; 2 [DB]; 3 [ETCD]")
)

type MysqlDB struct {
	Name            string
	Balancer        string
	DBName          string
	DBUser          string
	DBPass          string
	DBVariables     string
	ConnString      string
	ConnMaxIdleTime int
	ConnTimeout     int
	ConnMaxCnt      int
	ConnMaxLifetime int
	ConnEncoding    string
}

type MysqlSubSQL struct {
	SQL      string
	NoQuote  map[string]string
	NoEscape map[string]string
	Check    map[string]string
	Sharding map[string]string
	DB       string
}

type MysqlSQL struct {
	SQL      string
	NoQuote  map[string]string
	NoEscape map[string]string
	Check    map[string]string
	Sharding map[string]string
	DB       string
	SQLGroup []MysqlSubSQL
}

type MysqlCheck struct {
	Field string
	Regex string
}

type MysqlDBs map[string]MysqlDB
type MysqlSQLs map[string]MysqlSQL

type MysqlConfig struct {
	DC MysqlDBs
	SC MysqlSQLs
}

// loadConfigure loads configure from a JSON file Directory.
func (s *MysqlConfig) LoadConfigure() error {
	switch *mysqlLoadType {
	case 1:

		if err := loadFromFile(*mysqlDBPath, &s.DC); err != nil {
			return err
		}

		if err := loadFromFile(*mysqlSqlPath, &s.SC); err != nil {
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
		return errors.New("Error: Mysql Config Load Type")
	}
	return nil
}
