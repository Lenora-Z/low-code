package utils

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/tealeg/xlsx"
	"io/ioutil"
	"log"
	"os"
)

var excelDir = "./resource/excel"

var contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"

//创建Excel
func CreateExcel(fileName string, title []string, info []map[string]interface{}) (error, []byte, string, map[string]string) {
	//创建文件
	os.Mkdir(excelDir, os.ModePerm)
	f := xlsx.NewFile()
	sh, err := f.AddSheet("sheet1")
	if err != nil {
		log.Println("sheet not exist")
		return err, nil, "", nil
	}
	//表头
	row := sh.Row(0)
	row.WriteSlice(&title, -1)
	//内容
	for i, val := range info {
		row = sh.Row(i + 1)
		list := make([]interface{}, 0)
		for _, k := range title {
			if val[k] == nil {
				list = append(list, "")
			} else {
				list = append(list, val[k])
			}
		}
		row.WriteSlice(&list, -1)
	}
	path := fmt.Sprintf("%s/%s.xlsx", excelDir, fileName)
	logrus.Info(path)
	err = f.Save(path)
	if err != nil {
		log.Println("sheet not exist")
		return err, nil, "", nil
	}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		logrus.Error(err)
		return err, nil, "", nil
	}

	//delete file
	if err := os.Remove(path); err != nil {
		logrus.Error(err)
		return err, nil, "", nil
	}
	fileName = fileName + ".xlsx"
	extraHeaders := map[string]string{
		"Content-Disposition": `attachment; filename="` + fileName + `"`,
	}
	return nil, b, contentType, extraHeaders
}
