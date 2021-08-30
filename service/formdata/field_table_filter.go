package formdata

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

type FieldTableFilter struct {
	ID                         uint32 `gorm:"primaryKey;column:id;type:int(11) unsigned;not null" json:"id"`
	FormID                     uint32 `gorm:"column:form_id;type:int(11) unsigned;not null;default:0" json:"formId"`                                           // 表单id
	FieldID                    uint32 `gorm:"column:field_id;type:int(11) unsigned;not null;default:0" json:"fieldId"`                                         // 控件id
	DatasourceColumnRelationID uint32 `gorm:"column:datasource_column_relation_id;type:int(11) unsigned;not null;default:0" json:"datasourceColumnRelationId"` // 表字段关联表id,0表示本数据表字段
	DatasourceColumnID         uint32 `gorm:"column:datasource_column_id;type:int(11) unsigned;not null;default:0" json:"datasourceColumnId"`                  // 表字段id
	FieldType                  string `gorm:"column:field_type;type:varchar(20);not null" json:"fieldType"`                                                    // 控件类型
	FieldTypeCondition         int8   `gorm:"column:field_type_condition;type:tinyint(3);not null;default:1" json:"fieldTypeCondition"`                        // 控件筛选选项条件
	FieldTypeConditionValue    string `gorm:"column:field_type_condition_value;type:varchar(255);not null;default:''" json:"fieldTypeConditionValue"`          // 控件筛选选项值
}

// TableName get sql table name.获取数据库表名
func (f *FieldTableFilter) TableName() string {
	return "field_table_filter"
}

var fieldTableFilterTableName = (&FieldTableFilter{}).TableName()

func getFieldTableFilterByFieldId(db *gorm.DB, fieldId uint32) (bool, FieldTableFilterList) {
	var item FieldTableFilterList
	isNotFound := db.Table(fieldTableFilterTableName).Where("field_id = ?", fieldId).Find(&item).RecordNotFound()
	return isNotFound, item
}

func getWhereSentence(colName string, f *FieldTableFilter) []string {
	var whereSentences = make([]string, 0)
	var endNum int8
	if f.FieldType == "text" || f.FieldType == "radio" || f.FieldType == "date" || f.FieldType == "checkbox" {
		endNum = 5
	} else if f.FieldType == "num" {
		endNum = 7
	} else {
		return []string{}
	}
	if f.FieldTypeCondition < endNum {
		if f.FieldType == "radio" || f.FieldType == "checkbox" {
			conditionValues := strings.Split(f.FieldTypeConditionValue, ",")
			if len(conditionValues) > 1 { //"%s = '%s'"      "%s like '%%%s%%'"
				var condition string
				if f.FieldTypeCondition == 1 { //是
					condition = " = "
				} else if f.FieldTypeCondition == 2 { //不是
					condition = " != "
				} else if f.FieldTypeCondition == 3 { //包含
					condition = " in ("
				} else if f.FieldTypeCondition == 4 { //不包含
					condition = " not in ("
				}
				for _, value := range conditionValues {
					if f.FieldTypeCondition == 3 || f.FieldTypeCondition == 4 {
						condition += fmt.Sprintf("%s", value) + ","
					} else {
						sentence := colName + condition + fmt.Sprintf("%s", value)
						whereSentences = append(whereSentences, sentence)
					}
				}
				if len(condition) > 5 {
					whereSentences = []string{colName + strings.TrimRight(condition, ",") + ")"}
				}
			} else {
				whereSentences = []string{fmt.Sprintf(fieldCondition[f.FieldType][f.FieldTypeCondition], colName, f.FieldTypeConditionValue)}
			}
		} else if f.FieldType == "date" {
			//日期类型按天比较
			timestamp, _ := strconv.ParseInt(f.FieldTypeConditionValue, 10, 64)
			if f.FieldTypeCondition == 2 {
				timestamp += 3600 * 24
			}
			whereSentences = []string{fmt.Sprintf(fieldCondition[f.FieldType][f.FieldTypeCondition], colName, time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")[:10])}
		} else if f.FieldType == "num" {
			value, _ := strconv.Atoi(f.FieldTypeConditionValue)
			whereSentences = []string{fmt.Sprintf(fieldCondition[f.FieldType][f.FieldTypeCondition], colName, value)}
		} else {
			whereSentences = []string{fmt.Sprintf(fieldCondition[f.FieldType][f.FieldTypeCondition], colName, f.FieldTypeConditionValue)}
		}
	} else {
		whereSentences = []string{fmt.Sprintf(fieldCondition[f.FieldType][f.FieldTypeCondition], colName)}
	}

	logrus.Debug("whereSentences = ", whereSentences)
	return whereSentences
}
