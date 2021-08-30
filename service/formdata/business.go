package formdata

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"strings"
)

//获取列表数据
func getTableData(db *gorm.DB, sql string, page, limit uint32) (error, int, [][]string) {
	var count int
	var data = make([][]string, 0)

	sqlCount := "select count(*) " + sql[strings.Index(sql, " from "):]
	rows1, err := db.DB().Query(sqlCount)
	if err != nil {
		return err, 0, nil
	}
	defer rows1.Close()
	for rows1.Next() {
		_ = rows1.Scan(&count)
	}

	sqlValues := sql + fmt.Sprintf(" limit %d,%d", (page-1)*limit, limit)
	//logrus.Info(sqlValues)
	rows2, err := db.DB().Query(sqlValues)
	if err != nil {
		return err, 0, nil
	}
	defer rows2.Close()

	colTypes, _ := rows2.ColumnTypes()
	cols, _ := rows2.Columns()
	rawResult := make([][]byte, len(cols))
	dest := make([]interface{}, len(cols))
	for i, _ := range rawResult {
		dest[i] = &rawResult[i]
	}
	for rows2.Next() {
		result := make([]string, len(cols))
		_ = rows2.Scan(dest...)
		for i, raw := range rawResult {
			if raw == nil {
				if colTypes[i].DatabaseTypeName() == "DATETIME" {
					result[i] = "null"
				} else {
					result[i] = ""
				}
			} else {
				result[i] = string(raw)
			}
		}
		data = append(data, result)
	}

	return nil, count, data
}

//获取列表单条数据
func getTableDataById(db *gorm.DB, sql string) (error, []string) {
	rows, err := db.DB().Query(sql)
	if err != nil {
		return err, nil
	}

	colTypes, _ := rows.ColumnTypes()
	cols, _ := rows.Columns()
	rawResult := make([][]byte, len(cols))
	result := make([]string, len(cols))

	dest := make([]interface{}, len(cols))
	for i, _ := range rawResult {
		dest[i] = &rawResult[i]
	}
	for rows.Next() {
		_ = rows.Scan(dest...)
		for i, raw := range rawResult {
			if raw == nil {
				if colTypes[i].DatabaseTypeName() == "DATETIME" {
					result[i] = "null"
				} else {
					result[i] = ""
				}
			} else {
				result[i] = string(raw)
			}
		}
	}
	return nil, result
}

//更新某条列表数据
func updateTableData(db *gorm.DB, sql string) error {
	err := db.Exec(sql).Error
	return err
}

//删除某条列表数据
func deleteTableData(db *gorm.DB, sql string) error {
	err := db.Exec(sql).Error
	return err
}

//写入业务表数据
func insertData(db *gorm.DB, sql string) error {
	err := db.Exec(sql).Error
	return err
}

//获取某个字段的所有值
func getColValues(db *gorm.DB, sql string) (error, [][]string) {
	var data = make([][]string, 0)

	rows, err := db.DB().Query(sql)
	if err != nil {
		return err, nil
	}
	cols, _ := rows.Columns()
	rawResult := make([][]byte, len(cols))

	dest := make([]interface{}, len(cols))
	for i, _ := range rawResult {
		dest[i] = &rawResult[i]
	}
	for rows.Next() {
		result := make([]string, len(cols))
		_ = rows.Scan(dest...)
		for i, raw := range rawResult {
			if raw == nil {
				result[i] = ""
			} else {
				result[i] = string(raw)
			}
		}
		data = append(data, result)
	}

	return nil, data
}
