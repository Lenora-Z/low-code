// Package source
// Created by GoLand
// @User: lenora
// @Date: 2021/7/23
// @Time: 11:23

package source

var DatasourceTableColumns = struct {
	ID           string
	CreatedAt    string
	AppID        string
	SchemataID   string
	TableSchema  string
	TableName    string
	TableComment string
}{
	ID:           "id",
	CreatedAt:    "created_at",
	AppID:        "app_id",
	SchemataID:   "schemata_id",
	TableSchema:  "table_schema",
	TableName:    "table_name",
	TableComment: "table_comment",
}

var DatasourceColumnColumns = struct {
	ID                     string
	AppID                  string
	SchemataID             string
	TableSchema            string
	TableName              string
	DataType               string
	CharacterMaximumLength string
	ColumnName             string
	ColumnComment          string
	ShowType               string
	FieldType              string
	IsSystemField          string
}{
	ID:                     "id",
	AppID:                  "app_id",
	SchemataID:             "schemata_id",
	TableSchema:            "table_schema",
	TableName:              "table_name",
	DataType:               "data_type",
	CharacterMaximumLength: "character_maximum_length",
	ColumnName:             "column_name",
	ColumnComment:          "column_comment",
	ShowType:               "show_type",
	FieldType:              "field_type",
	IsSystemField:          "is_system_field",
}

var DatasourceColumnRelationColumns = struct {
	ID             string
	SourceTableID  string
	SourceColumnID string
	TargetTableID  string
	TargetColumnID string
	Type           string
}{
	ID:             "id",
	SourceTableID:  "source_table_id",
	SourceColumnID: "source_column_id",
	TargetTableID:  "target_table_id",
	TargetColumnID: "target_column_id",
	Type:           "type",
}

var DatasourceMetadataColumns = struct {
	ID                 string
	AppID              string
	DatasourceColumnID string
	Group              string
	Key                string
	Value              string
}{
	ID:                 "id",
	AppID:              "app_id",
	DatasourceColumnID: "datasource_column_id",
	Group:              "group",
	Key:                "key",
	Value:              "value",
}

const (
	TEXT         = "text"     //文本
	NUM          = "num"      //数字
	DATE         = "date"     //日期
	RADIO        = "radio"    //单选
	CHECKBOX     = "checkbox" //多选
	USER         = "user"     //成员
	ORGANIZATION = "org"      //组织
	RELATION     = "relation" //关联
	FILE         = "file"     //文件
)
