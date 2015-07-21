package main

import (
	"database/sql"
	"fmt"
	"github.com/Unknwon/goconfig"
	"github.com/alecthomas/log4go"
	_ "github.com/go-sql-driver/mysql"
	"path/filepath"
	"time"
)

type Env_t struct {
	Conf *goconfig.ConfigFile
	Db   *sql.DB
	Log  log4go.Logger
}

var Env *Env_t

func Load(path, conf string) error {
	Env = &Env_t{}
	var err error

	Env.Conf, err = goconfig.LoadConfigFile(filepath.Join(path, conf))
	if err != nil {
		return fmt.Errorf("loading config %s error", filepath.Join(path, conf))
	}

	logpath, err := Env.Conf.GetValue("log", "path")
	if err != nil {
		return fmt.Errorf("loading log:path error")
	}

	logfile, err := Env.Conf.GetValue("log", "file")
	if err != nil {
		return fmt.Errorf("loading log:file error")
	}

	Env.Log = make(log4go.Logger)

	Env.Log.AddFilter("info", log4go.INFO, log4go.NewFileLogWriter(filepath.Join(logpath, logfile), false))
	Env.Log.AddFilter("err", log4go.WARNING, log4go.NewFileLogWriter(filepath.Join(logpath, logfile+".wf"), false))

	schema, err := Env.Conf.GetValue("mysql", "schema")
	if err != nil {
		return fmt.Errorf("loading mysql:schema error")
	}

	Env.Db, err = sql.Open("mysql", schema)
	if err != nil {
		return fmt.Errorf("open mysql connection error, reason : %s", err.Error())
	}

	return nil
}

func Close() {
	Env.Log.Close()
	time.Sleep(100) //log4go has channel flush bug, sleep a while can fix this problem
}

func (*Env_t) GetDB() *sql.DB {
	return Env.Db
}
