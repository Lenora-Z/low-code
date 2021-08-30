package formdata

import (
	"github.com/jinzhu/gorm"
)

type DatasourceColumn struct {
	ID                     uint32 `gorm:"primaryKey;column:id;type:int(11) unsigned;not null" json:"id"`
	AppID                  uint32 `gorm:"column:app_id;type:int(11) unsigned;not null;default:0" json:"appId"`                     // 应用id
	SchemataID             uint32 `gorm:"column:schemata_id;type:int(11) unsigned;not null;default:0" json:"schemataId"`           // 数据源id
	TableSchema            string `gorm:"column:table_schema;type:varchar(64);not null;default:''" json:"tableSchema"`             // 数据库名称
	TableName              string `gorm:"column:table_name;type:varchar(64);not null;default:''" json:"tableName"`                 // 数据表名称
	DataType               string `gorm:"column:data_type;type:varchar(64);not null;default:''" json:"dataType"`                   // 字段数据类型
	CharacterMaximumLength uint64 `gorm:"column:character_maximum_length;type:bigint(21) unsigned" json:"characterMaximumLength"`  // 字段长度
	ColumnName             string `gorm:"column:column_name;type:varchar(64);not null;default:''" json:"columnName"`               // 字段名称
	ColumnComment          string `gorm:"column:column_comment;type:varchar(1024);not null;default:''" json:"columnComment"`       // 字段备注
	ShowType               string `gorm:"column:show_type;type:varchar(20);not null" json:"showType"`                              // 字段显示类型(业务类型) text:文本，num: 数字 date: 日期 radio:单选 checkbox:多选 user:成员 org:组织 relation: 关联 file: 文件
	FieldType              string `gorm:"column:field_type;type:varchar(20);not null" json:"fieldType"`                            // 控件类型
	IsSystemField          bool   `gorm:"column:is_system_field;type:tinyint(1) unsigned;not null;default:0" json:"isSystemField"` // 是否为系统字段 0:不是系统字段 1:是系统字段
}

// GetTableName get sql table name.获取数据库表名
func (f *DatasourceColumn) GetTableName() string {
	return "datasource_column"
}

var datasourceColumnTableName = (&DatasourceColumn{}).GetTableName()

func getDatasourceColumnById(db *gorm.DB, id uint32) (bool, *DatasourceColumn) {
	var item DatasourceColumn
	isNotFound := db.Table(datasourceColumnTableName).Where("id = ?", id).First(&item).RecordNotFound()
	return isNotFound, &item
}

func getBusinessColNameByIds(db *gorm.DB, ids []uint32) (error, []string, []string, []uint32) {
	var datasourceColumnList DatasourceColumnList
	err := db.Table(datasourceColumnTableName).Where("id in (?)", ids).Order("id asc").Find(&datasourceColumnList).Error

	colNames := datasourceColumnList.ColNames()
	colViewNames := datasourceColumnList.ColViewNames()
	choiceIds := datasourceColumnList.ChoiceIds() //需要映射的字段

	return err, colNames, colViewNames, choiceIds
}

//isSystemField=0，查询非系统字段；isSystemField=1，查询系统字段；isSystemField=3，查询所有字段
func getBusinessColNameByTableName(db *gorm.DB, appId uint32, tableName string, isSystemField uint8) (error, []string, []uint32, []string, []string) {
	var datasourceColumnList DatasourceColumnList
	query := db.Table(datasourceColumnTableName).Where("app_id = ?", appId).Where("table_name = ?", tableName)
	if isSystemField == 0 {
		query = query.Where("is_system_field = ?", 0)
	} else if isSystemField == 1 {
		query = query.Where("is_system_field = ?", 1)
	}

	err := query.Order("id asc").Find(&datasourceColumnList).Error

	colNames := datasourceColumnList.ColNames()
	colIds := datasourceColumnList.ColIds()
	colViewNames := datasourceColumnList.ColViewNames()
	colFieldTypes := datasourceColumnList.ColFieldTypes()

	return err, colNames, colIds, colViewNames, colFieldTypes
}

func isCreatorIdColumn(db *gorm.DB, colId uint32) bool {
	var item DatasourceColumn
	isNotFound := db.Table(datasourceColumnTableName).Where("id = ?", colId).Where("column_name = ?", "creator_id").First(&item).RecordNotFound()
	return isNotFound
}
