// Package data_db
// Created by GoLand
// @User: lenora
// @Date: 2021/7/31
// @Time: 16:18

package data_db

type ConditionItem struct {
	ColumnName string      `json:"column_name"` //字段名称
	Cond       uint8       `json:"cond"`        //条件 1-是/等于  2-否/不等于
	Value      interface{} `json:"value"`       //条件值
}

type ConditionGroup []ConditionItem
