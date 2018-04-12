package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc/grpclog"

	conf "github.com/toniz/gudp/config"
	pb "github.com/toniz/gudp/interface"
)

type DBHandle map[string]*sql.DB
type MysqlDataService struct {
	dbh DBHandle
	sc  conf.MysqlSQLs
}

type MysqlConfigLoader interface {
	LoadConfigure() error
}

// NewMysqlDataService creates a MysqlDataService with the provided GrpcDataServer.
// Call loadConfigure To Load Configure Data
func NewMysqlDataService() (*MysqlDataService, error) {
	s := &MysqlDataService{dbh: make(DBHandle), sc: make(conf.MysqlSQLs)}
	l, err := conf.NewMysqlConfig()
	if err != nil || len(l.DC) == 0 {
		return nil, err
	}

	s.sc = l.SC
	for k, v := range l.DC {
		// Construct Connect String
		dbstr := v.DBUser + `:` + v.DBPass + `@tcp(` + v.ConnString + `)/` + v.DBName + `?charset=` + v.ConnEncoding + v.DBVariables
		grpclog.Infof("Loading Config Connect To: %s", dbstr)

		if s.dbh[k], err = MysqlConnect(dbstr); err != nil {
			grpclog.Errorf("DB[%s] Connect Failed [%s]: %v", k, dbstr, err)
		}

		// Set ConnMaxCnt
		if v.ConnMaxCnt != 0 {
			s.dbh[k].SetMaxOpenConns(v.ConnMaxCnt)
		}

		// Set ConnMaxLifetime
		if v.ConnMaxLifetime != 0 {
			duration := strconv.Itoa(v.ConnMaxLifetime)
			if d, e := time.ParseDuration(duration + "s"); e == nil {
				s.dbh[k].SetConnMaxLifetime(d)
			}
		}
	}
	return s, err
}

// MysqlConnect connect to mysql, and retrun sql.DB.
// conn : root:bbwhat@tcp(127.0.0.1:4000)/bbwhat
func MysqlConnect(conn string) (db *sql.DB, err error) {
	// Open database connection
	db, err = sql.Open("mysql", conn)
	if err != nil {
		return
	}

	err = db.Ping()
	return
}

// MysqlRealEscapeString. Golang don`t have mysql_real_escape_string function.
// Using this funtion to make escape string.
func MysqlRealEscapeString(value string) string {
	value = strings.Replace(value, `\`, `\\`, -1)
	value = strings.Replace(value, `"`, `\"`, -1)
	return value
}

// MysqlQuery Using db Execute the sqlstr.
func MysqlQuery(db *sql.DB, sqlstr string, res *pb.Response) error {
	rows, err := db.Query(sqlstr)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	// Make a slice for the values
	values := make([]string, len(columns))

	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// Fetch rows
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return err
		}

		// Now do something with the data.
		// Here we just print each column as a string.
		data := make(map[string]string)
		for i, col := range values {
			data[columns[i]] = col
		}

		if rows.Err() == nil {
			var row pb.RowData
			row.Fields = data
			res.Rows = append(res.Rows, &row)
		} else {
			err = rows.Err()
		}
	}

	return err
}

func (s *MysqlDataService) DoCommit(ctx context.Context, req *pb.Query) (*pb.Response, error) {

	if len(req.Ident) < 1 {
		return nil, errors.New("req.Ident Must Define.")
	}

	switch req.Opt {
	case "1":
		return s.autoCommit(ctx, req)
	case "2":
		return s.transCommit(ctx, req)
	case "3":
		return s.multiInsert(ctx, req)
	default:
		if len(s.sc[req.Ident].SQLGroup) > 0 {
			return s.transCommit(ctx, req)
		} else {
			return s.autoCommit(ctx, req)
		}
	}
	return nil, nil
}

func (s *MysqlDataService) autoCommit(ctx context.Context, req *pb.Query) (*pb.Response, error) {

	var res pb.Response
	var err error

	sqlc := s.sc[req.Ident].SQL
	if len(sqlc) < 1 {
		msg := "Error: Not Found This Sql ID On AutoCommit Mode." + req.Ident
		grpclog.Warningf(msg)
		return &res, errors.New(msg)
	}

	dbname := s.sc[req.Ident].DB

	for k, v := range req.Params {
		// Check the parameter using regex
		// eg: "check":   {"id": "^\\d+$"}
		// The id parameter must be number string.
		if rex, ok := s.sc[req.Ident].Check[k]; ok {
			var validParam = regexp.MustCompile(rex)
			if match := validParam.MatchString(v); !match {
				grpclog.Warningf("Param Check Failed regex[%s] param[%s]", rex, v)
				msg := "Parameter Check Failed!"
				return &res, errors.New(msg)
			}
		}

		// Escape query string
		// eg: "noescape":{"values":""}
		// The `values`parameter don`t need Escape
		var val string
		if _, ok := s.sc[req.Ident].NoEscape[k]; ok {
			val = v
		} else {
			val = MysqlRealEscapeString(v)
		}

		// Replace parameter
		if _, ok := s.sc[req.Ident].NoQuote[k]; ok {
			sqlc = strings.Replace(sqlc, "$"+k+"$", val, -1)
		} else {
			sqlc = strings.Replace(sqlc, "$"+k+"$", "\""+val+"\"", -1)
		}

		// DB Sharding Support
		if _, ok := s.sc[req.Ident].Sharding[k]; ok {
			dbname = strings.Replace(dbname, "$"+k+"$", val, -1)
		}
	}

	grpclog.Infof("Sql[%s] Dbname[%s]", sqlc, dbname)

	if dbh := s.dbh[dbname]; dbh != nil {
		err = MysqlQuery(dbh, sqlc, &res)
		if err != nil {
			grpclog.Warningf(" %v", err)
		}
	} else {
		msg := fmt.Sprintf("Error: [%s]Not Found In Configure..", dbname)
		err = errors.New(msg)
		grpclog.Errorf(msg)
	}
	return &res, err
}

func (s *MysqlDataService) transCommit(ctx context.Context, req *pb.Query) (*pb.Response, error) {

	var res pb.Response
	var err error

	var dbs []string
	var sqls []string
	for i, sg := range s.sc[req.Ident].SQLGroup {
		sqlc := sg.SQL
		if len(sqlc) < 1 {
			msg := fmt.Sprintf("Error: [%s] Not Found This[%d] Sql On Transaction Commit Mode.", req.Ident, i)
			grpclog.Warningf(msg)
			return nil, errors.New(msg)
		}

		dbname := sg.DB
		if len(dbname) < 1 {
			msg := fmt.Sprintf("Error: [%s]Not Found This[%d] DBName On Transaction Commit Mode.", req.Ident, i)
			grpclog.Warningf(msg)
			return nil, errors.New(msg)
		}

		for k, v := range req.Group[i].Params {
			// Check the parameter using regex
			// eg: "check":   {"id": "^\\d+$"}
			// The id parameter must be number string.
			if rex, ok := sg.Check[k]; ok {
				var validParam = regexp.MustCompile(rex)
				if match := validParam.MatchString(v); !match {
					msg := fmt.Sprintf("sql[%d] Param Check Failed regex[%s] param[%s]", i, rex, v)
					grpclog.Warningf(msg)
					return nil, errors.New(msg)
				}
			}

			// Escape query string
			// eg: "noescape":{"values":""}
			// The `values`parameter don`t need Escape
			var val string
			if _, ok := sg.NoEscape[k]; ok {
				val = v
			} else {
				val = MysqlRealEscapeString(v)
			}

			// Replace parameter
			if _, ok := sg.NoQuote[k]; ok {
				sqlc = strings.Replace(sqlc, "$"+k+"$", val, -1)
			} else {
				sqlc = strings.Replace(sqlc, "$"+k+"$", "\""+val+"\"", -1)
			}

			// DB Sharding Support
			if _, ok := sg.Sharding[k]; ok {
				dbname = strings.Replace(dbname, "$"+k+"$", val, -1)
			}
		}
		grpclog.Infof("Sql[%s] ", sqlc)
		grpclog.Infof("Dbname[%s] ", dbname)
		dbs = append(dbs, dbname)
		sqls = append(sqls, sqlc)
	}

	// Get DB Handle, Then Set AutoCommit = false
	rollback := false
	dbhs := make(map[string]*sql.DB)
	for _, dbname := range dbs {
		// Unique
		if _, ok := dbhs[dbname]; !ok {
			dbh := s.dbh[dbname]
			if dbh == nil {
				msg := "Error: [" + req.Ident + "]No Such DB Handle." + dbname
				err = errors.New(msg)
				grpclog.Warningf(msg)
				rollback = true
				break
			}

			if _, err = dbh.ExecContext(ctx, "SET AUTOCOMMIT=0;"); err != nil {
				msg := "Warn: Set AutoCommit=0 Failed: " + dbname
				err = errors.New(msg)
				grpclog.Warningf(msg)
				rollback = true
				break
			}
			dbhs[dbname] = dbh
		}
	}

	// Exec mysql query
	if !rollback {
		for i, q := range sqls {
			grpclog.Infof("Exec [%d] Sql[%s] ", i, q)

			err = MysqlQuery(dbhs[dbs[i]], q, &res)
			if err != nil {
				grpclog.Warningf("MysqlQuery Failed: seq[%d] sql[%s] db[%s]: %v", i, q, dbs[i], err)
				rollback = true
				break
			}
		}
	}

	for _, dbh := range dbhs {
		if rollback {
			if _, e := dbh.ExecContext(ctx, "ROLLBACK;"); e != nil {
				grpclog.Warningf("Warn: Rollback Transcation Failed: ", req.Ident, e)
			}
		} else {
			if _, err = dbh.ExecContext(ctx, "COMMIT;"); err != nil {
				grpclog.Warningf("Warn: Commit Transcation Failed: %s [%v]", req.Ident, err)
			}
		}

		if _, e := dbh.ExecContext(ctx, "SET AUTOCOMMIT=1;"); e != nil {
			msg := "Warn: Set AutoCommit=1 Failed"
			grpclog.Warningf(msg)
		}
	}

	return &res, err
}

func (s *MysqlDataService) multiInsert(ctx context.Context, req *pb.Query) (*pb.Response, error) {
	var res pb.Response
	var err error

	return &res, err
}
