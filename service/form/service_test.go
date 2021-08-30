// Package form
// Created by GoLand
// @User: lenora
// @Date: 2021/8/18
// @Time: 17:45

package form

import (
	"testing"
)

func TestFieldButtons_Fields(t *testing.T) {
	e := FormNumber([]byte("{\"colspan\":24,\"description\":\"只能输入数字\",\"digit\":\"fhajsk\",\"explain\":\"\",\"funcCategory\":[\"collection\"],\"id\":\"item-Number-6674ED9C65ED682\",\"limitAmount\":false,\"max\":10,\"min\":null,\"name\":\"Number6674ED9C65ED682\",\"placeholder\":\"请输入数值\",\"pluginCategory\":\"input\",\"repeat\":false,\"require\":false,\"sourceFieldId\":1,\"title\":\"应聘者id\",\"type\":\"Number\",\"typeName\":\"数值\",\"unit\":\"元\"}"))
	t.Log(e)
}
