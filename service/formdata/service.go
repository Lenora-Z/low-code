package formdata

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

type FormdataService interface {
	//通过field_id获取表field_table数据
	GetFieldTableByFieldId(fieldId uint32) (bool, *FieldTable)
	//通过field_id获取表field_table_column数据
	GetFieldTableColumnByFieldId(fieldId uint32) (bool, FieldTableColumnList)
	//通过field_id获取表field_table_filter数据
	GetFieldTableFilterByFieldId(fieldId uint32) (bool, FieldTableFilterList)
	//通过id获取表datasource_column数据
	GetDatasourceColumnById(id uint32) (bool, *DatasourceColumn)
	//通过id获取表datasource_column_relation数据
	GetDatasourceColumnRelationById(id uint32) (bool, *DatasourceColumnRelation)
	//通过id获取表datasource_table数据
	GetDatasourceTableById(id uint32) (bool, *DatasourceTable)
	//解析成条件sql
	GetWhereSentence(fieldTableFilter *FieldTableFilter) (error, []string)
	//更新列表单条数据
	UpdateTableData(businessDb *gorm.DB, tableName string, rowId uint32, newValues string) error
	//删除列表单条数据
	DeleteTableData(businessDb *gorm.DB, tableName string, rowId uint32) error
	//获取业务数据表表字段  ids需升序
	GetBusinessColNameByIds(ids []uint32) (error, []string, []string, []uint32)
	//获取业务数据表的数据
	GetBusinessTableData(businessDb *gorm.DB, mainTables, relationTables, filterCondition map[string][]string, page, limit uint32) (error, int, [][]string, []uint32)
	//获取业务数据表单条的数据详情
	GetBusinessTableDataById(businessDb *gorm.DB, mainTables, relationTables map[string][]string, id uint32) (error, []string)
	//通过表名table_name获取业务数据表表字段
	GetBusinessColNameByTableName(appId uint32, tableName string) (error, []string, []uint32, []string, []string)
	//根据id判断是否为创建人字段
	IsCreatorIdColumn(db *gorm.DB, colId uint32) bool
	//写入业务表数据
	InsertBusinessData(businessDb *gorm.DB, tableName string, colData map[string][]interface{}) error
	//获取某个字段的所有值
	GetColValues(businessDb *gorm.DB, tableName, colName, whereSentence string) (error, [][]string)
	//通过field_id获取表field_records数据
	GetFieldRecordsByFieldId(fieldId uint32) (bool, *FieldRecords)
	//通过datasource_column_id获取表datasource_metadata数据
	GetDatasourceMetadataById(datasourceColumnId uint32) (bool, []DatasourceMetadata)
	//通过field_id获取表field_button数据
	GetFieldButtonByFieldId(fieldId uint32) (bool, *FieldButton)
	//通过field_id获取表field_linkage数据
	GetFieldLinkageByFieldId(fieldId uint32) (bool, *FieldLinkage)
	//通过key获取映射值
	GetDatasourceMetadataByKey(id uint32, key string) (bool, string)
	//通过value获取key
	GetDatasourceMetadataByValue(id uint32, value string) (bool, string)
	//通过field_id获取表field_table_button数据
	GetFieldTableButtonByFieldId(fieldId uint64, name string) (bool, *FieldTableButton)
}

type formdataService struct {
	db *gorm.DB
}

func NewFormdataService(db *gorm.DB) FormdataService {
	u := new(formdataService)
	u.db = db
	return u
}

//通过field_id获取表field_table数据
func (srv *formdataService) GetFieldTableByFieldId(fieldId uint32) (bool, *FieldTable) {
	return getFieldTableByFieldId(srv.db, fieldId)
}

//通过field_id获取表field_table_column数据
func (srv *formdataService) GetFieldTableColumnByFieldId(fieldId uint32) (bool, FieldTableColumnList) {
	return getFieldTableColumnByFieldId(srv.db, fieldId)
}

//通过field_id获取表field_table_filter数据
func (srv *formdataService) GetFieldTableFilterByFieldId(fieldId uint32) (bool, FieldTableFilterList) {
	return getFieldTableFilterByFieldId(srv.db, fieldId)
}

//通过id获取表datasource_table数据
func (srv *formdataService) GetDatasourceTableById(id uint32) (bool, *DatasourceTable) {
	return getDatasourceTableById(srv.db, id)
}

//通过id获取表datasource_column数据
func (srv *formdataService) GetDatasourceColumnById(id uint32) (bool, *DatasourceColumn) {
	return getDatasourceColumnById(srv.db, id)
}

//通过id获取表datasource_column_relation数据
func (srv *formdataService) GetDatasourceColumnRelationById(id uint32) (bool, *DatasourceColumnRelation) {
	return getDatasourceColumnRelationById(srv.db, id)
}

//解析成条件sql
func (srv *formdataService) GetWhereSentence(fieldTableFilter *FieldTableFilter) (error, []string) {
	isNotFound, datasourceColumn := getDatasourceColumnById(srv.db, fieldTableFilter.DatasourceColumnID)
	if isNotFound {
		return errors.New("not found"), nil
	}
	return nil, getWhereSentence(datasourceColumn.ColumnName, fieldTableFilter)
}

//更新列表单条数据
func (srv *formdataService) UpdateTableData(businessDb *gorm.DB, tableName string, rowId uint32, newValues string) error {
	sql := "update " + tableName + " set " + newValues + fmt.Sprintf(" where id = %d", rowId)
	return updateTableData(businessDb, sql)
}

//删除列表单条数据
func (srv *formdataService) DeleteTableData(businessDb *gorm.DB, tableName string, rowId uint32) error {
	sql := "delete from " + tableName + fmt.Sprintf(" where id = %d", rowId)
	return deleteTableData(businessDb, sql)
}

//通过表名ids获取业务数据表表字段  ids需升序
func (srv *formdataService) GetBusinessColNameByIds(ids []uint32) (error, []string, []string, []uint32) {
	return getBusinessColNameByIds(srv.db, ids)
}

//通过表名table_name获取业务数据表表字段
func (srv *formdataService) GetBusinessColNameByTableName(appId uint32, tableName string) (error, []string, []uint32, []string, []string) {
	return getBusinessColNameByTableName(srv.db, appId, tableName, 0) //isSystemField=0，查询非系统字段；isSystemField=1，查询系统字段；isSystemField=3，查询所有字段
}

//根据id判断是否为创建人字段
func (srv *formdataService) IsCreatorIdColumn(db *gorm.DB, colId uint32) bool {
	return isCreatorIdColumn(db, colId)
}

//获取业务数据表的数据
func (srv *formdataService) GetBusinessTableData(businessDb *gorm.DB, mainTables, relationTables, filterCondition map[string][]string, page, limit uint32) (error, int, [][]string, []uint32) {
	var mainTableName string
	var relateColIds = make([]uint32, 0)
	var selSql, showCols, leftJoin, whereSentences string
	if len(mainTables) != 0 {
		for tName, colNames := range mainTables {
			mainTableName = tName
			for _, colName := range colNames {
				showCols += tName + "." + colName + ","
			}
		}
	} else {
		logrus.Error("there is no main table")
		return errors.New("there is no main table"), 0, nil, nil
	}

	if len(relationTables) != 0 {
		for tName, colNames := range relationTables {
			//因为map类型是无序的，所以需要在这里记录字段的id顺序
			for _, v := range strings.Split(colNames[2], ",") {
				vUint64, _ := strconv.ParseUint(v, 10, 32)
				relateColIds = append(relateColIds, uint32(vUint64))
			}
			for i := 3; i < len(colNames); i++ {
				showCols += tName + "." + colNames[i] + ","
			}
			leftJoin += " left join " + tName + " on " + mainTableName + "." + colNames[0] + " = " + tName + "." + colNames[1]
		}
	} else {
		logrus.Debug("there is no relation table")
	}
	showCols = mainTableName + ".Id," + strings.TrimRight(showCols, ",")

	if len(filterCondition) != 0 {
		for tName, filters := range filterCondition {
			for _, filter := range filters {
				if strings.Contains(filter, "length(") {
					whereSentences += filter[:7] + tName + "." + filter[7:] + " and "
				} else {
					whereSentences += tName + "." + filter + " and "
				}
			}
		}
		whereSentences = " where " + strings.TrimRight(whereSentences, "and ")
	} else {
		logrus.Debug("there is no filter condition")
	}

	order := " order by " + mainTableName + ".updated_at desc"
	selSql = "select " + showCols + " from " + mainTableName + leftJoin + whereSentences + order
	err, count, data := getTableData(businessDb, selSql, page, limit)
	if err != nil {
		return err, 0, nil, nil
	}

	return nil, count, data, relateColIds
}

//获取业务数据表单条的数据详情
func (srv *formdataService) GetBusinessTableDataById(businessDb *gorm.DB, mainTables, relationTables map[string][]string, id uint32) (error, []string) {
	var mainTableName string
	var selSql, showCols, leftJoin, whereSentences string
	if len(mainTables) != 0 {
		for tName, colNames := range mainTables {
			mainTableName = tName
			for _, colName := range colNames {
				showCols += tName + "." + colName + ","
			}
		}
	} else {
		logrus.Error("there is no main table")
		return errors.New("there is no main table"), nil
	}

	if len(relationTables) != 0 {
		for tName, colNames := range relationTables {
			for i := 2; i < len(colNames); i++ {
				showCols += tName + "." + colNames[i] + ","
			}
			leftJoin += " left join " + tName + " on " + mainTableName + "." + colNames[0] + " = " + tName + "." + colNames[1]
		}
	} else {
		logrus.Debug("there is no relation table")
	}
	showCols = mainTableName + ".Id," + strings.TrimRight(showCols, ",")

	whereSentences = " where " + mainTableName + fmt.Sprintf(".id = %d", id)
	selSql = "select " + showCols + " from " + mainTableName + leftJoin + whereSentences
	err, data := getTableDataById(businessDb, selSql)
	if err != nil {
		return err, nil
	}

	return nil, data
}

//写入业务表数据
func (srv *formdataService) InsertBusinessData(businessDb *gorm.DB, tableName string, colData map[string][]interface{}) error {
	var items string
	for colName, value := range colData {
		if value[0].(string) == "int" {
			items += colName + " = " + fmt.Sprintf("%d", value[1]) + ","
		} else if value[0].(string) == "datetime" {
			items += colName + " = " + fmt.Sprintf("'%s'", time.Unix(value[1].(int64), 0).Format("2006-01-02 15:04:05")) + ","
		} else {
			items += colName + " = " + fmt.Sprintf("'%s'", value[1]) + ","
		}
	}

	sql := "insert into " + tableName + " set " + strings.TrimRight(items, ",")

	return insertData(businessDb, sql)
}

//通过field_id获取表field_records数据
func (srv *formdataService) GetFieldRecordsByFieldId(fieldId uint32) (bool, *FieldRecords) {
	return getFieldRecordsByFieldId(srv.db, fieldId)
}

//通过datasource_column_id获取表datasource_metadata数据
func (srv *formdataService) GetDatasourceMetadataById(datasourceColumnId uint32) (bool, []DatasourceMetadata) {
	return getDatasourceMetadataById(srv.db, datasourceColumnId)
}

//通过field_id获取表field_button数据
func (srv *formdataService) GetFieldButtonByFieldId(fieldId uint32) (bool, *FieldButton) {
	return getFieldButtonByFieldId(srv.db, fieldId)
}

//通过field_id获取表field_linkage数据
func (srv *formdataService) GetFieldLinkageByFieldId(fieldId uint32) (bool, *FieldLinkage) {
	return getFieldLinkageByFieldId(srv.db, fieldId)
}

//获取某个字段的所有值
func (srv *formdataService) GetColValues(businessDb *gorm.DB, tableName, colName, whereSentence string) (error, [][]string) {
	if whereSentence != "" {
		whereSentence = " where " + whereSentence
	}
	sql := fmt.Sprintf("select id,%s from %s %s order by %s", colName, tableName, whereSentence, colName)
	return getColValues(businessDb, sql)
}

//通过key获取映射值
func (srv *formdataService) GetDatasourceMetadataByKey(id uint32, key string) (bool, string) {
	return getDatasourceMetadataByKey(srv.db, id, key)
}

//通过value获取key
func (srv *formdataService) GetDatasourceMetadataByValue(id uint32, value string) (bool, string) {
	return getDatasourceMetadataByValue(srv.db, id, value)
}

//通过field_id获取表field_table_button数据
func (srv *formdataService) GetFieldTableButtonByFieldId(fieldId uint64, name string) (bool, *FieldTableButton) {
	return getFieldTableButtonByFieldId(srv.db, fieldId, name)
}
