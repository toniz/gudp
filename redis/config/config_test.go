package config

import (
	"testing"
)

func Test_ParseCfgFromDir(t *testing.T) {
	dbpath := "../../conf/db/"
	sqlpath := "../../conf/sql/"

	dbcfg := make(DBCfgs)
	sqlcfg := make(SQLCfgs)
	if err := loadFromFile(dbpath, &dbcfg); err != nil {
		t.Error("Parse DB Config From Path Failed", err)
	}

	if err := loadFromFile(sqlpath, &sqlcfg); err != nil {
		t.Error("Parse SQL Config From Path Failed", err)
	}
}
