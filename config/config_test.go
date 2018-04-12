package config

import (
	"testing"
)

func Test_ParseCfgFromDir(t *testing.T) {
	dbpath := "../conf/mysql/db/"
	sqlpath := "../conf/mysql/sql/"
	redispath := "../conf/redis/srv/"

	dbcfg := make(DBCfgs)
	sqlcfg := make(SQLCfgs)
	rdcfg := make(RedisSRVs)
	if err := loadFromFile(dbpath, &dbcfg); err != nil {
		t.Error("Parse DB Config From Path Failed", err)
	}

	if err := loadFromFile(sqlpath, &sqlcfg); err != nil {
		t.Error("Parse SQL Config From Path Failed", err)
	}

	if err := loadFromFile(redispath, &rdcfg); err != nil {
		t.Error("Parse Redis Config From Path Failed", err)
	}
}
