// Package source
// Created by GoLand
// @User: lenora
// @Date: 2021/7/22
// @Time: 10:41

package source

import (
	"github.com/Lenora-Z/low-code/utils"
	"github.com/jinzhu/gorm"
)

type Service interface {
	// GetAllTable 获取应用内全部数据表
	GetAllTable(appId uint64) (error, uint32, TableList)
	// GetColumnList 数据表字段列表
	GetColumnList(appId uint64, tbName string) (error, []DatasourceColumn)
	// GetColumnListByDataType 获取字段列表by业务类型
	GetColumnListByDataType(appId uint64, name string, types ...string) (error, []DatasourceColumn)
	// GetColumnListByIds 数据表字段列表by多个id
	GetColumnListByIds(id ...uint64) (error, []DatasourceColumn)
	// GetColumnListByTbNames 数据字段列表by多个数据表名
	GetColumnListByTbNames(appId uint64, name ...string) (error, []DatasourceColumn)
	// GetMetaListByColumnId 获取字典信息by字段id
	GetMetaListByColumnId(id uint64) (error, []DatasourceMetadata)
	// GetTableDetail 获取数据表
	GetTableDetail(id uint64) (bool, *DatasourceTable)
	// GetTableByName 获取数据表by表名称
	GetTableByName(name string, appId uint64) (bool, *DatasourceTable)
	// GetTableList 数据表列表
	GetTableList(page, limit uint32, appId uint64, name string) (error, uint32, TableList)
	// GetTableListByIds 获取数据表列表by多个id
	GetTableListByIds(appId uint64, id ...uint64) (error, uint32, TableList)
	// GetTableRelations 获取数据表关联关系
	GetTableRelations(id uint64) (error, RelationList)
}

type service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) Service {
	u := new(service)
	u.db = db
	return u
}

func (srv *service) GetColumnList(appId uint64, tbName string) (error, []DatasourceColumn) {
	return getColumnList(srv.db, appId, []string{tbName}, nil, nil)
}

func (srv *service) GetColumnListByDataType(appId uint64, name string, types ...string) (error, []DatasourceColumn) {
	return getColumnList(srv.db, appId, []string{name}, nil, types)
}

func (srv *service) GetColumnListByIds(id ...uint64) (error, []DatasourceColumn) {
	return getColumnList(srv.db, 0, nil, id, nil)
}

func (srv *service) GetColumnListByTbNames(appId uint64, name ...string) (error, []DatasourceColumn) {
	return getColumnList(srv.db, appId, name, nil, nil)
}

func (srv *service) GetAllTable(appId uint64) (error, uint32, TableList) {
	return getTableList(srv.db, 0, utils.MAX_LIMIT, appId, nil, "")
}

func (srv *service) GetTableDetail(id uint64) (bool, *DatasourceTable) {
	return tableDetail(srv.db, id)
}

func (srv *service) GetTableList(page, limit uint32, appId uint64, name string) (error, uint32, TableList) {
	offset := (page - 1) * limit
	return getTableList(srv.db, offset, limit, appId, nil, name)
}

func (srv *service) GetTableListByIds(appId uint64, id ...uint64) (error, uint32, TableList) {
	return getTableList(srv.db, 0, utils.MAX_LIMIT, appId, id, "")
}

func (srv *service) GetTableByName(name string, appId uint64) (bool, *DatasourceTable) {
	column := make(map[string]interface{})
	if appId != 0 {
		column[DatasourceTableColumns.AppID] = appId
	}
	column[DatasourceTableColumns.TableName] = name
	return tableDetailByColumn(srv.db, column)
}

func (srv *service) GetTableRelations(id uint64) (error, RelationList) {
	return getRelationList(srv.db, id)
}

func (srv *service) GetMetaListByColumnId(id uint64) (error, []DatasourceMetadata) {
	return getMetaList(srv.db, id)
}
