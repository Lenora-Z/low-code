// Package data_db
// Created by GoLand
// @User: lenora
// @Date: 2021/7/31
// @Time: 16:00

package data_db

import (
	"fmt"
	"github.com/gohouse/gorose"
	"github.com/sirupsen/logrus"
)

type Service interface {
	NewItem(col map[string]interface{}) error
	UpdateData(col map[string]interface{}, cond []ConditionGroup) error
	GetListWithGroup(group string) (error, []map[string]interface{})
	GetTotal(cond map[string]interface{}) (error, int64)
	GetTotalBetweenCond(field string, start, end interface{}, cond map[string]interface{}) (error, int64)
	DeleteItem(cond map[string]interface{}) error
	GetList(cond map[string]interface{}) (error, []map[string]interface{})
	GetItem(id uint64) (error, map[string]interface{})

	ContractFiling(id, sum uint64) error
}

type service struct {
	db      *gorose.Connection
	name    string
	session *gorose.Session
}

func NewService(db *gorose.Connection, name string) Service {
	u := new(service)
	u.db = db
	u.name = name
	return u
}

func NewTxService(db *gorose.Session, name string) Service {
	u := new(service)
	u.name = name
	u.session = db
	u.db = db.Connection
	return u
}

func (srv *service) NewItem(col map[string]interface{}) error {
	_, err := srv.table().Data(col).Insert()
	return err
}

func (srv *service) GetItem(id uint64) (error, map[string]interface{}) {
	res, err := srv.table().Where("id", id).First()
	return err, res
}

func (srv *service) UpdateData(col map[string]interface{}, cond []ConditionGroup) error {
	query := srv.table().Data(col)
	if len(cond) > 0 {
		for i, x := range cond {
			var f string
			for d, t := range x {
				symbol := "="
				if t.Cond == 2 {
					symbol = "!="
				}
				if d == 0 {
					f = fmt.Sprintf(`%s %s '%v'`, t.ColumnName, symbol, t.Value)
				} else {
					f = f + fmt.Sprintf(` and %s %s '%v'`, t.ColumnName, symbol, t.Value)
				}

			}
			logrus.Info("where:", f)
			if i == 0 {
				query = query.Where(f)
			} else {
				query = query.OrWhere(f)
			}
		}
	} else {
		query = query.Where("id", ">", 0)
	}
	c, err := query.Update()
	logrus.Info("affected rows:", c)
	return err
}

func (srv *service) GetList(cond map[string]interface{}) (error, []map[string]interface{}) {
	query := srv.table()
	if len(cond) <= 0 {
		query = query.Where("1 = 1")
	} else {
		for k, v := range cond {
			query = query.Where(k, v)
		}
	}
	res, err := query.Get()
	return err, res
}

func (srv *service) GetListWithGroup(group string) (error, []map[string]interface{}) {
	query := srv.table()
	res, err := query.Fields(group + ",count(*) as count").Group(group).Get()

	return err, res
}

func (srv *service) GetTotal(cond map[string]interface{}) (error, int64) {
	query := srv.table()
	if len(cond) <= 0 {
		query = query.Where("1 = 1")
	} else {
		for k, v := range cond {
			query = query.Where(k, v)
		}
	}
	count, err := query.Count()
	return err, count
}

func (srv *service) GetTotalBetweenCond(field string, start, end interface{}, cond map[string]interface{}) (error, int64) {
	query := srv.table().WhereBetween(field, []interface{}{start, end})
	for k, v := range cond {
		query = query.Where(k, v)
	}

	count, err := query.Count()
	return err, count

}

func (srv *service) DeleteItem(cond map[string]interface{}) error {
	query := srv.table()
	for k, v := range cond {
		query = query.Where(k, v)
	}

	count, err := query.Delete()
	logrus.Info("affected rows:", count)
	return err
}

func (srv *service) ContractFiling(id uint64, sum uint64) error {
	srv.name = "hrm_labor_contract"
	data := map[string]interface{}{
		"sum":         sum,
		"is_received": 1,
	}
	res, err := srv.table().Where("id", id).Data(data).Update()
	logrus.Info("affected rows:", res)
	return err
}

func (srv *service) table() *gorose.Session {
	var query *gorose.Session
	if srv.session != nil {
		query = srv.session.Table(srv.name)
	} else {
		query = srv.db.Table(srv.name)
	}
	return query
}
