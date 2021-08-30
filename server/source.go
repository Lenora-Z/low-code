// Package server
// Created by GoLand
// @User: lenora
// @Date: 2021/7/22
// @Time: 10:42

package server

import (
	"github.com/Lenora-Z/low-code/service/source"
	"github.com/Lenora-Z/low-code/utils"
	"github.com/gin-gonic/gin"
)

type TableListPageVO struct {
	List  []source.DatasourceTable `json:"list"`
	Count uint32                   `json:"count"`
}

// TableListByPage
// @Summary 数据表列表(分页)[new]
// @Tags table
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @param page query string true "页码"
// @param limit query string true "页容量"
// @param name query string false "表名/描述"
// @Success 200 {object} ApiResponse{result=TableListPageVO}
// @Router /table/page [get]
func (ds *defaultServer) TableListByPage(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*CustomClaims)
	pageStr := ctx.Query("page")
	limitStr := ctx.Query("limit")
	name := ctx.Query("name")

	if pageStr == "" || limitStr == "" {
		ds.InvalidParametersError(ctx)
		return
	}

	srv := source.NewService(ds.db)
	err, count, list := srv.GetTableList(utils.NewStr(pageStr).Uint32(), utils.NewStr(limitStr).Uint32(), claims.AppId, name)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	ds.ResponseSuccess(ctx, TableListPageVO{
		List:  list,
		Count: count,
	})
	return
}

type TableAttr struct {
	Id      uint64 `json:"id"`
	Name    string `json:"name"`
	Comment string `json:"comment"`
	Schema  string `json:"schema"`
}

type TableColumnVO struct {
	Item    *TableAttr                `json:"item"`    //表详情
	Columns []source.DatasourceColumn `json:"columns"` //字段列表
}

// TableColumns
// @Summary 数据表字段列表[new]
// @Tags table
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @param name query string true "表名"
// @param type query string true "0-仅获取字段列表 1-字段列表+表详情"
// @Success 200 {object} ApiResponse{result=TableColumnVO}
// @Router /table/columns [get]
func (ds *defaultServer) TableColumns(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*CustomClaims)
	name := ctx.Query("name")
	tStr := ctx.Query("type")
	if name == "" || tStr == "" {
		ds.InvalidParametersError(ctx)
		return
	}

	t := utils.NewStr(tStr).Uint()
	srv := source.NewService(ds.db)
	var item TableAttr
	if t == 1 {
		notFound, tb := srv.GetTableByName(name, claims.AppId)
		if notFound {
			ds.ResponseError(ctx, NOT_FOUND, "source table not found")
			return
		}
		item = TableAttr{
			Id:      tb.ID,
			Name:    tb.TableName,
			Comment: tb.TableComment,
			Schema:  tb.TableSchema,
		}
	}

	err, list := srv.GetColumnList(claims.AppId, name)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	ds.ResponseSuccess(ctx, TableColumnVO{
		Item:    &item,
		Columns: list,
	})
	return
}

type SourceRelationItem struct {
	SourceColumnName   string `json:"source_column_name"`   //源字段名称
	TargetTableName    string `json:"target_table_name"`    //目标表名称
	TargetTableComment string `json:"target_table_comment"` //目标表描述
	TargetColumnName   string `json:"target_column_name"`   //目标字段名称
	Type               uint8  `json:"type"`                 //关联类型.1-1对1，2-1对多,3-多对1,4-多对多
}

type SourceRelationVO struct {
	Item     TableAttr            `json:"item"`
	Relation []SourceRelationItem `json:"relation"`
}

// TableRelations
// @Summary 数据表关联关系(数据表查看关联用)[new]
// @Tags table
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @param id query string true "表id"
// @Success 200 {object} ApiResponse{result=SourceRelationVO}
// @Router /table/relations [get]
func (ds *defaultServer) TableRelations(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*CustomClaims)
	idStr := ctx.Query("id")
	if idStr == "" {
		ds.InvalidParametersError(ctx)
		return
	}

	id := utils.NewStr(idStr).Uint64()
	srv := source.NewService(ds.db)
	//获取数据表详情
	notFound, item := srv.GetTableDetail(id)
	if notFound || item.AppID != claims.AppId {
		ds.ResponseError(ctx, NOT_FOUND, "source table not found")
		return
	}

	//获取关联关系
	err, relations := srv.GetTableRelations(id)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	sourceColumnsIds := relations.SourceColumns()
	targetColumnsIds := relations.TargetColumns()
	columnIds := make([]uint64, 0, cap(sourceColumnsIds)+cap(targetColumnsIds))
	columnIds = append(sourceColumnsIds, targetColumnsIds...)

	err, _, targets := srv.GetTableListByIds(claims.AppId, relations.TargetTables()...)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	err, columns := srv.GetColumnListByIds(columnIds...)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	tableRelation := make([]SourceRelationItem, 0, cap(relations))
	for _, x := range relations {
		var i SourceRelationItem
		i.Type = x.Type
		for _, t := range targets {
			if t.ID == x.TargetTableID {
				i.TargetTableName = t.TableName
				i.TargetTableComment = t.TableComment
				break
			}
		}
		for _, sc := range columns {
			if sc.ID == x.SourceColumnID {
				i.SourceColumnName = sc.Name
				break
			}
		}
		for _, tc := range columns {
			if tc.ID == x.TargetColumnID {
				i.TargetColumnName = tc.Name
				break
			}
		}
		tableRelation = append(tableRelation, i)
	}

	ds.ResponseSuccess(ctx, SourceRelationVO{
		Item: TableAttr{
			Id:      item.ID,
			Name:    item.TableName,
			Comment: item.TableComment,
			Schema:  item.TableSchema,
		},
		Relation: tableRelation,
	})
	return
}

type SourceRelationItemAttrVO struct {
	Id              uint64 `json:"id"`
	TargetTableID   uint64 `json:"target_table_id"`   //目标表id
	TargetTableName string `json:"target_table_name"` //目标表名称
}

// TableRelationsAttr
// @Summary 数据表关联关系(关联记录用)[new]
// @Tags table
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @param id query string true "表id"
// @Success 200 {object} ApiResponse{result=[]SourceRelationItemAttrVO}
// @Router /table/relations/attr  [get]
func (ds *defaultServer) TableRelationsAttr(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*CustomClaims)
	idStr := ctx.Query("id")
	if idStr == "" {
		ds.InvalidParametersError(ctx)
		return
	}

	id := utils.NewStr(idStr).Uint64()
	srv := source.NewService(ds.db)

	//获取关联关系
	err, relations := srv.GetTableRelations(id)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	err, _, targets := srv.GetTableListByIds(claims.AppId, relations.TargetTables()...)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	err, columns := srv.GetColumnListByIds(relations.TargetColumns()...)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	tableRelation := make([]SourceRelationItemAttrVO, 0, cap(relations))
	for _, x := range relations {
		var i = SourceRelationItemAttrVO{Id: x.ID}
		for _, t := range targets {
			if t.ID == x.TargetTableID {
				i.TargetTableName = t.TableName
				break
			}
		}
		for _, sc := range columns {
			if sc.ID == x.SourceColumnID {
				i.TargetTableName = sc.Name
				break
			}
		}
		tableRelation = append(tableRelation, i)
	}
	ds.ResponseSuccess(ctx, tableRelation)
	return
}

type TableModelVO struct {
	RelationId uint64                    `json:"relation_id"` //关联id 本表值为0
	Id         uint64                    `json:"id"`          //表id
	Name       string                    `json:"name"`        //表名称
	Comment    string                    `json:"comment"`     //表描述
	Columns    []source.DatasourceColumn `json:"columns"`     //数据列
}

// TableModel
// @Summary 数据表模型/本表及字表的所有数据[new]
// @Tags table
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @param id query string true "表id"
// @Success 200 {object} ApiResponse{result=[]TableModelVO}
// @Router /table/model  [get]
func (ds *defaultServer) TableModel(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*CustomClaims)
	idStr := ctx.Query("id")
	if idStr == "" {
		ds.InvalidParametersError(ctx)
		return
	}

	id := utils.NewStr(idStr).Uint64()
	srv := source.NewService(ds.db)
	//获取数据表详情
	notFound, item := srv.GetTableDetail(id)
	if notFound || item.AppID != claims.AppId {
		ds.ResponseError(ctx, NOT_FOUND, "source table not found")
		return
	}

	err, relations := srv.GetTableRelations(id)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	err, _, targets := srv.GetTableListByIds(claims.AppId, relations.TargetTables()...)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	tables := make([]string, 0, cap(targets)*2)
	tables = append(tables, item.TableName)
	tables = append(tables, targets.Names()...)

	err, columns := srv.GetColumnListByTbNames(claims.AppId, tables...)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	tableModel := make([]TableModelVO, 0, cap(tables))
	itemCols := make([]source.DatasourceColumn, 0, cap(columns))
	for _, cols := range columns {
		if cols.TableSchema == item.TableSchema && cols.TableName == item.TableName {
			itemCols = append(itemCols, cols)
		}
	}
	tableModel = append(tableModel, TableModelVO{
		RelationId: 0,
		Id:         item.ID,
		Name:       item.TableName,
		Comment:    item.TableComment,
		Columns:    itemCols,
	})
	for _, x := range relations {
		for _, t := range targets {
			if t.ID == x.TargetTableID {
				cls := make([]source.DatasourceColumn, 0, cap(columns))
				for _, c := range columns {
					if c.TableName == t.TableName && c.SchemataID == t.SchemataID {
						cls = append(cls, c)
						continue
					}
				}
				tableModel = append(tableModel, TableModelVO{
					RelationId: x.ID,
					Id:         t.ID,
					Name:       t.TableName,
					Comment:    t.TableComment,
					Columns:    cls,
				})
				break
			}
		}
	}

	ds.ResponseSuccess(ctx, tableModel)
	return
}

// TableList
// @Summary 数据表完整列表[new]
// @Tags table
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @Success 200 {object} ApiResponse{result=[]source.DatasourceTable}
// @Router /table/list  [get]
func (ds *defaultServer) TableList(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*CustomClaims)
	srv := source.NewService(ds.db)
	err, _, list := srv.GetAllTable(claims.AppId)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	ds.ResponseSuccess(ctx, list)
	return
}

// TableRelationColumns
// @Summary 关联字段列表[new]
// @Tags table
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @param id query string true "表id"
// @Success 200 {object} ApiResponse{result=[]TableModelVO}
// @Router /table/relation/columns  [get]
func (ds *defaultServer) TableRelationColumns(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*CustomClaims)
	idStr := ctx.Query("id")
	if idStr == "" {
		ds.InvalidParametersError(ctx)
		return
	}

	id := utils.NewStr(idStr).Uint64()
	srv := source.NewService(ds.db)

	err, relations := srv.GetTableRelations(id)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	err, _, targets := srv.GetTableListByIds(claims.AppId, relations.TargetTables()...)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	err, columns := srv.GetColumnListByTbNames(claims.AppId, targets.Names()...)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	tableModel := make([]TableModelVO, 0, cap(relations))
	for _, x := range relations {
		for _, t := range targets {
			if t.ID == x.TargetTableID {
				cls := make([]source.DatasourceColumn, 0, cap(columns))
				for _, c := range columns {
					if c.TableName == t.TableName && c.SchemataID == t.SchemataID {
						cls = append(cls, c)
						continue
					}
				}
				tableModel = append(tableModel, TableModelVO{
					RelationId: x.ID,
					Id:         t.ID,
					Name:       t.TableName,
					Comment:    t.TableComment,
					Columns:    cls,
				})
				break
			}
		}
	}

	ds.ResponseSuccess(ctx, tableModel)
	return
}

type MetaListVO struct {
	ColumnId uint64 `json:"column_id"` //字段id
	Group    string `json:"group"`     //组名
	Key      string `json:"key"`       //键名
	Value    string `json:"value"`     //键值
}

// TypeColumns
// @Summary 数据表的字段字段列表(分类获取)[new]
// @Tags table
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @param name query string true "表名称"
// @param type query string true "类型 1-获取可进行筛选的字段 2-时间类型字段"
// @Success 200 {object} ApiResponse{result=[]source.DatasourceColumn}
// @Router /table/columns/type  [get]
func (ds *defaultServer) TypeColumns(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*CustomClaims)
	name := ctx.Query("name")
	tStr := ctx.Query("type")
	if name == "" || tStr == "" {
		ds.InvalidParametersError(ctx)
		return
	}

	t := utils.NewStr(tStr).Uint8()

	srv := source.NewService(ds.db)
	var err error
	var list []source.DatasourceColumn
	switch t {
	case 1:
		err, list = srv.GetColumnListByDataType(claims.AppId, name, source.TEXT, source.NUM, source.DATE, source.RADIO, source.CHECKBOX)
	case 2:
		err, list = srv.GetColumnListByDataType(claims.AppId, name, source.DATE)
	}

	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	ds.ResponseSuccess(ctx, list)
	return

}

// TableMeta
// @Summary 数据表字段字典内容[new]
// @Tags table
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @param id query string true "字段id"
// @Success 200 {object} ApiResponse{result=[]MetaListVO}
// @Router /table/meta  [get]
func (ds *defaultServer) TableMeta(ctx *gin.Context) {
	idStr := ctx.Query("id")
	if idStr == "" {
		ds.InvalidParametersError(ctx)
		return
	}

	id := utils.NewStr(idStr).Uint64()
	srv := source.NewService(ds.db)
	err, meta := srv.GetMetaListByColumnId(id)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	list := make([]MetaListVO, 0, cap(meta))
	for _, x := range meta {
		list = append(list, MetaListVO{
			ColumnId: x.DatasourceColumnID,
			Group:    x.Group,
			Key:      x.Key,
			Value:    x.Value,
		})
	}

	ds.ResponseSuccess(ctx, list)
	return
}
