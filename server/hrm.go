// Package server
// Created by GoLand
// @User: lenora
// @Date: 2021/8/13
// @Time: 10:38

package server

import (
	"encoding/json"
	"github.com/Lenora-Z/low-code/service/data_db"
	"github.com/Lenora-Z/low-code/utils"
	"github.com/gin-gonic/gin"
	"github.com/gohouse/gorose"
	"github.com/sirupsen/logrus"
	"time"
)

type MonthVO struct {
	Month     int64 `json:"month"`         //月份
	Entry     int64 `json:"entry_num"`     //入职
	Dimission int64 `json:"dimission_num"` //离职
}

type StatisticVO struct {
	Total     int64           `json:"total"`      //总人数
	ForEntry  int64           `json:"for_entry"`  //待入职
	Informal  int64           `json:"informal"`   //待转正
	ToBeLeft  int64           `json:"to_be_left"` //待离职
	Male      int64           `json:"male"`       //男
	Female    int64           `json:"female"`     //女
	Education map[int64]int64 `json:"education"`  //key值  1:博士;2:硕士;3:本科;4:大专;5:大专以下
	Month     []MonthVO       `json:"month"`      //各月份
}

// GetHrmStatistics
// @Summary 获取hrm统计数据
// @Tags hrm
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @Success 200 {object} ApiResponse{result=StatisticVO}
// @Router /hrm/statistic [get]
func (ds *defaultServer) GetHrmStatistics(ctx *gin.Context) {
	employeeSrv := data_db.NewService(ds.dataDb, "hrm_employee")

	var item StatisticVO

	err, total := employeeSrv.GetTotal(nil)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	item.Total = total

	err, statusList := employeeSrv.GetListWithGroup("status")
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	for _, x := range statusList {
		s, ok := x["status"].(int64)
		c, cok := x["count"].(int64)
		if !ok || !cok {
			continue
		}
		switch s {
		case 4:
			item.Informal = c
		case 5:
			item.ToBeLeft = c
		}
	}

	err, genderList := employeeSrv.GetListWithGroup("gender")
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	for _, x := range genderList {
		s, ok := x["gender"].(int64)
		c, cok := x["count"].(int64)
		if !ok || !cok {
			continue
		}
		switch s {
		case 1:
			item.Male = c
		case 2:
			item.Female = c
		}
	}

	err, eduList := employeeSrv.GetListWithGroup("education")
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	edu := make(map[int64]int64, 0)
	for _, x := range eduList {
		s, ok := x["education"].(int64)
		c, cok := x["count"].(int64)
		if !ok || !cok {
			continue
		}
		edu[s] = c
	}
	item.Education = edu

	applicantSrv := data_db.NewService(ds.dataDb, "hrm_applicant")
	err, entry := applicantSrv.GetTotal(map[string]interface{}{"status": 6})
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	item.ForEntry = entry

	statistic := data_db.NewService(ds.dataDb, "hrm_employee_statistic")
	err, log := statistic.GetList(map[string]interface{}{
		"year": 2021,
	})
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	jByte, err := json.Marshal(log)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	logList := make([]MonthVO, 0, cap(log))
	if err := json.Unmarshal(jByte, &logList); err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	nowYear, nowMonth, _ := time.Now().Date()
	//nowYear := 2021
	//nowMonth := time.Month(7)
	date := time.Date(nowYear, nowMonth, 1, 0, 0, 0, 0, time.Local)
	start := date.Format(utils.DateFormatStr)
	endTime := date.AddDate(0, 1, 0)
	end := endTime.Format(utils.DateFormatStr)

	err, entry, dimission := getLastMonth(ds.dataDb, start, end)
	logList = append(logList, MonthVO{
		Month:     int64(nowMonth),
		Entry:     entry,
		Dimission: dimission,
	})
	item.Month = logList

	ds.ResponseSuccess(ctx, item)
	return

}

func (cs *cronServer) SetMonthData() {
	nowYear, nowMonth, _ := time.Now().Date()
	//nowYear := 2021
	//nowMonth := time.Month(7)
	date := time.Date(nowYear, nowMonth, 1, 0, 0, 0, 0, time.Local)
	end := date.Format(utils.DateFormatStr)
	startTime := date.AddDate(0, -1, 0)
	year, month, _ := startTime.Date()
	start := startTime.Format(utils.DateFormatStr)
	err, entry, dimission := getLastMonth(cs.dataDb, start, end)
	if err != nil {
		logrus.Error(err)
		return
	}

	db := cs.dataDb.NewSession()
	db.Begin()
	srv := data_db.NewTxService(db, "hrm_employee_statistic")
	if err := srv.DeleteItem(map[string]interface{}{
		"year":  int64(year),
		"month": int64(month),
	}); err != nil {
		db.Rollback()
		logrus.Error(err)
		return
	}

	if err := srv.NewItem(map[string]interface{}{
		"year":          int64(year),
		"month":         int64(month),
		"entry_num":     entry,
		"dimission_num": dimission,
	}); err != nil {
		db.Rollback()
		logrus.Error(err)
		return
	}

	db.Commit()
	return

}

func getLastMonth(db *gorose.Connection, start, end string) (error, int64, int64) {

	err, entry := data_db.NewService(db, "hrm_employee").GetTotalBetweenCond("entry_time", start, end, nil)
	if err != nil {
		logrus.Error(err)
		return err, 0, 0
	}

	err, dimission := data_db.NewService(db, "hrm_resignation").GetTotalBetweenCond("resign_time", start, end, nil)

	if err != nil {
		logrus.Error(err)
		return err, 0, 0
	}

	return nil, entry, dimission
}
