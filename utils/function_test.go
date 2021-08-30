//Created by Goland
//@User: lenora
//@Date: 2021/3/10
//@Time: 7:18 下午
package utils

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestPassword(t *testing.T) {
	str, err := CryptoPassword("calendar")
	logrus.Info(str, "---", err)
	if err != nil {
		logrus.Error("error:", err)
	}

	status := CheckPassword(str, "calendar")
	logrus.Info(status)
}

func TestGetUUid(t *testing.T) {
	for i := 0; i < 8; i++ {
		t.Log(GetUUid())
	}
}

type TestStruct struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestCreateExcel(t *testing.T) {
	title := []string{"a", "b", "c"}
	info := []map[string]interface{}{
		{"a": 1, "b": 2, "c": nil},
		{"a": 3, "b": nil, "c": 23},
		{"a": 5, "b": 6, "c": nil},
	}
	infoByte, _ := json.Marshal(info)
	t.Log(string(infoByte))
	err, src, contentType, header := CreateExcel("test", title, info)
	if err != nil {
		t.Fatal(err.Error())
	}
	logrus.Info(err, src, contentType, header)
}
