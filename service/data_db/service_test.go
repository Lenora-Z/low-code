// Package data_db
// Created by GoLand
// @User: lenora
// @Date: 2021/8/13
// @Time: 10:54

package data_db

import (
	"fmt"
	"github.com/Lenora-Z/low-code/utils"
	"github.com/gohouse/gorose"
	"github.com/sirupsen/logrus"
	"testing"
	"time"
)

var config = struct {
	User     string `default:"root" yaml:"user"`
	Password string `default:"" yaml:"password"`
	Name     string `yaml:"ip"`
	Port     uint   `default:"3306" yaml:"port"`
	DbName   string `required:"true" yaml:"db_name"`
	Charset  string `default:"utf8" yaml:"charset"`
	MaxIdle  int    `default:"10" yaml:"max_idle"`
	MaxOpen  int    `default:"50" yaml:"max_open"`
	LogMode  bool   `yaml:"log_mode"`
	Loc      string `required:"true" yaml:"loc"`
}{
	"root", "123456", "192.168.3.47", 3306, "bpmn_app", "utf8", 10, 100, true, " Asia/Shanghai",
}

var db *gorose.Connection

func getDataDbConnection() (*gorose.Connection, error) {
	conf := config
	format := "%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True"
	dsn := fmt.Sprintf(format, conf.User, conf.Password, conf.Name, conf.Port, conf.DbName, conf.Charset)
	conn, err := gorose.Open(&gorose.DbConfigSingle{
		Driver:          "mysql",
		EnableQueryLog:  conf.LogMode,
		SetMaxOpenConns: conf.MaxOpen,
		SetMaxIdleConns: conf.MaxIdle,
		Prefix:          "",
		Dsn:             dsn,
	})
	if err != nil {
		return nil, err
	}
	conn.Use(gorose.NewLogger())
	return conn, nil
}

func init() {
	conn, err := getDataDbConnection()
	if err != nil {
		logrus.Error(err)
		return
	} else {
		db = conn
	}

}

func TestService_GetListWithGroup(t *testing.T) {
	NewService(db, "hrm_employee").GetListWithGroup("status")
}

func TestService_GetTotalBetweenCond(t *testing.T) {
	year, month, _ := time.Now().Date()
	date := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	start := date.AddDate(0, -1, 0)
	srv := NewService(db, "hrm_employee")
	srv.GetTotalBetweenCond("entry_time", start.Format(utils.DateFormatStr), date.Format(utils.DateFormatStr), map[string]interface{}{"id": 1})
	srv.GetList(nil)
}

func TestService_ContractFiling(t *testing.T) {
	srv := NewService(db, "hrm_labor_contract")
	err, res := srv.GetItem(1)
	if err != nil {
		t.Log(err)
		return
	}

	sum, _ := res["sum"].(int64)
	if sum != 0 {
		sum = sum - 1
	}
	srv.ContractFiling(1, uint64(sum))
}
