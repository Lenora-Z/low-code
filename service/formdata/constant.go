package formdata

import (
	"sort"
)

var fieldCondition = map[string]map[int8]string{
	"text":     {1: "%s = '%s'", 2: "%s != '%s'", 3: "%s like '%%%s%%'", 4: "%s not like '%%%s%%'", 5: "length(%s) = 0", 6: "length(%s) != 0"},
	"num":      {1: "%s = %d", 2: "%s != %d", 3: "%s > %d", 4: "%s < %d", 5: "%s >= %d", 6: "%s <= %d", 7: "length(%s) = 0", 8: "length(%s) != 0"},
	"radio":    {1: "%s = %s", 2: "%s != %s", 3: "%s in (%s)", 4: "%s not in (%s)", 5: "%s = 0", 6: "%s != 0"},
	"checkbox": {1: "%s = %s", 2: "%s != %s", 3: "%s in (%s)", 4: "%s not in (%s)", 5: "%s = 0", 6: "%s != 0"},
	"date":     {1: "%s < '%s'", 2: "%s > '%s'", 3: "%s like '%s%%'", 4: "%s not like '%s%%'", 5: "length(%s) = 0", 6: "length(%s) != 0"},
}

type FieldTableColumnList []FieldTableColumn
type FieldTableFilterList []FieldTableFilter
type DatasourceColumnList []DatasourceColumn

func (d DatasourceColumnList) ColNames() []string {
	var colNames = make([]string, 0)
	for _, v := range d {
		colNames = append(colNames, v.ColumnName)
	}
	return colNames
}

func (d DatasourceColumnList) ColIds() []uint32 {
	var colIds = make([]uint32, 0)
	for _, v := range d {
		colIds = append(colIds, v.ID)
	}
	return colIds
}

func (d DatasourceColumnList) ColFieldTypes() []string {
	var colFieldTypes = make([]string, 0)
	for _, v := range d {
		colFieldTypes = append(colFieldTypes, v.FieldType)
	}
	return colFieldTypes
}

func (d DatasourceColumnList) ColViewNames() []string {
	var colViewNames = make([]string, 0)
	for _, v := range d {
		colViewNames = append(colViewNames, v.ColumnComment)
	}
	return colViewNames
}

//返回是单选、多选控件的字段id
func (d DatasourceColumnList) ChoiceIds() []uint32 {
	var choiceIds = make([]uint32, 0)
	for _, v := range d {
		if v.FieldType == "MultipleChoice" || v.FieldType == "SingleChoice" {
			choiceIds = append(choiceIds, v.ID)
		}
	}
	return choiceIds
}

func (f FieldTableColumnList) DatasourceColIds() []uint32 {
	var colIds = make([]uint32, 0)
	for _, v := range f {
		colIds = append(colIds, v.DatasourceColumnID)
	}
	//升序排序
	sort.Slice(colIds, func(i, j int) bool {
		return colIds[i] < colIds[j]
	})
	return colIds
}

func (f FieldTableColumnList) RelationIdsAndColIds() map[uint32][]uint32 {
	var relationIdsAndColIds = make(map[uint32][]uint32, 0)
	for _, v := range f {
		if _, ok := relationIdsAndColIds[v.DatasourceColumnRelationID]; !ok {
			relationIdsAndColIds[v.DatasourceColumnRelationID] = []uint32{v.DatasourceColumnID}
		} else {
			relationIdsAndColIds[v.DatasourceColumnRelationID] = append(relationIdsAndColIds[v.DatasourceColumnRelationID], v.DatasourceColumnID)
		}
	}
	//升序排序
	for k, _ := range relationIdsAndColIds {
		sort.Slice(relationIdsAndColIds[k], func(i, j int) bool {
			return relationIdsAndColIds[k][i] < relationIdsAndColIds[k][j]
		})
	}
	return relationIdsAndColIds
}
