// Package server
// Created by GoLand
// @User: lenora
// @Date: 2021/8/13
// @Time: 15:06

package server

import (
	"fmt"
	"github.com/Lenora-Z/low-code/conf"
	"github.com/gohouse/gorose"
	"github.com/jinzhu/configor"
	"github.com/jinzhu/gorm"
)

type cronServer struct {
	name   string
	task   string
	conf   *conf.ServerConfig
	db     *gorm.DB
	dataDb *gorose.Connection
}

func (cs *cronServer) Run(configPath string) error {
	// config
	if err := cs.config(configPath); err != nil {
		return fmt.Errorf("rs.config(): %s", err.Error())
	}

	// mysql
	if err := cs.dbClient(); err != nil {
		return fmt.Errorf("rs.dbClient(): %s", err.Error())
	}

	// init
	if err := cs.init(); err != nil {
		return fmt.Errorf("rs.init(): %s", err.Error())
	}

	// exec task
	return cs.taskExec()
}

func (cs *cronServer) Close() error {

	if cs.db != nil {
		if err := cs.db.Close(); err != nil {
			return err
		}
	}

	return nil
}

func (cs *cronServer) config(configPath string) error {
	cs.conf = new(conf.ServerConfig)
	err := configor.Load(cs.conf, configPath)
	if err != nil {
		return err
	}
	return nil
}

func (cs *cronServer) dbClient() error {
	db, err := getDbConnection(cs.conf.DbConfig, 0, 0)
	if err != nil {
		return fmt.Errorf("mysql connect error:%+v", err)
	}
	fmt.Println("mysql connect successfully")
	cs.db = db

	dataDb, err := getDataDbConnection(cs.conf.BusinessDbConfig, 0, 0)
	if err != nil {
		return fmt.Errorf("data db connect error:%+v", err)
	}
	fmt.Println("data db connect successfully")
	cs.dataDb = dataDb
	return nil
}

func (cs *cronServer) init() error {
	return nil
}

func (cs *cronServer) taskExec() error {
	switch cs.task {
	case "hrm_statistic":
		cs.SetMonthData()
		return nil
	}
	return nil

}
