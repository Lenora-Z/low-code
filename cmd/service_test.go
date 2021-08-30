package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Lenora-Z/low-code/service/form"
	"github.com/gohouse/gorose"
	"github.com/sirupsen/logrus"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"
	"unicode/utf8"
)

func Test(t *testing.T) {
	var text = "Table832131687970xhjkshddshda-kdasiowqnc"
	prep := regexp.MustCompile(`^\w{1,40}$`)
	t.Log(prep.MatchString(text))

	//字符串长度
	str := ""
	l1 := len([]rune(str))
	l2 := bytes.Count([]byte(str), nil) - 1
	l3 := strings.Count(str, "") - 1
	l4 := utf8.RuneCountInString(str)
	fmt.Println(l1)
	fmt.Println(l2)
	fmt.Println(l3)
	fmt.Println(l4)

	ss := "Process_A2DEEF9:1:c62960bd-c445-11eb-8a48-96d78af4641a"

	key := strings.Index(ss, ":")
	t.Log(key)
	t.Log(ss[0:key])

}

func StructToMapDemo(obj interface{}) map[string]interface{} {
	obj1 := reflect.TypeOf(obj)
	obj2 := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < obj1.NumField(); i++ {
		data[obj1.Field(i).Name] = obj2.Field(i).Interface()
	}
	return data
}
func TestStructToMap(t *testing.T) {
	student := struct {
		Code      uint32    `json:"code"`
		Sign      string    `json:"sign"`
		Age       uint32    `json:"age"`
		UpdatedAt time.Time `json:"updated_at"`
	}{10, "jqw", 18, time.Now()}
	aa, _ := json.Marshal(student)
	data := make(map[string]interface{})
	json.Unmarshal(aa, &data)
	fmt.Println(data)
}

type TreeData struct {
	Id       uint64 `json:"id"`
	Data     []map[string]interface{}
	Children []TreeData
}

func TestTree(t *testing.T) {
	tree := form.FormTreeList{
		form.FormTree{
			Id: 160,
			Children: form.FormTreeList{
				form.FormTree{
					Id: 188,
					Children: form.FormTreeList{
						form.FormTree{
							Id:       187,
							Children: form.FormTreeList{},
						},
						form.FormTree{
							Id:       147,
							Children: form.FormTreeList{},
						},
					},
				},
			},
		},
	}
	aa := TreeData{}
	aa.set(tree)
	t.Log(aa)
}

func goRose() {
	//conf := ds.conf.DbConfig
	//format := "%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=%s"
	//dsn := fmt.Sprintf(format, conf.User, conf.Password, conf.Name, conf.Port, conf.DbName, conf.Charset, url.QueryEscape(conf.Loc))
	dsn := ""
	conn, err := gorose.Open(&gorose.DbConfigSingle{
		Driver:          "mysql",
		EnableQueryLog:  true,
		SetMaxOpenConns: 0,
		SetMaxIdleConns: 0,
		Prefix:          "",
		Dsn:             dsn,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	db := conn.NewSession()
	item := map[string]interface{}{
		"name":       "aaa",
		"age":        3,
		"updated_at": time.Now(),
	}

	a, err := db.Table("tt").Data(item).InsertGetId()
	logrus.Info(a, err)
}

func (node *TreeData) set(list form.FormTreeList) {
	for _, x := range list {
		data := TreeData{
			Id:       x.Id,
			Data:     []map[string]interface{}{{"aa": "bb"}, {"cc": "dd"}},
			Children: []TreeData{},
		}
		data.set(x.Children)
		node.Children = append(node.Children, data)
	}
	return
}

func TestRegex(t *testing.T) {
	str := `恭喜您{aaa}被我们公司录用, 请于 {1,2} 到我们公司{2}报道{564}`
	reg := regexp.MustCompile(`{(\d|,)+}`)
	m := reg.FindAllString(str, -1)
	for i, x := range m {
		m[i] = x[1 : len(x)-1]
	}
}

type NodeCommon struct {
	Key   uint64      `json:"key"`         // 节点唯一标识
	Name  string      `json:"name"`        // 节点名称
	Desc  string      `json:"description"` // 节点描述
	Type  uint        `json:"type"`        // 节点类型
	Props interface{} `json:"props"`       //节点属性
	List  []Node      `json:"list"`
}

type Node []NodeCommon

func TestDecode(t *testing.T) {
	str := "[{\"key\":5894158697,\"name\":\"开始节点\",\"description\":\"节点描述\",\"type\":0,\"props\":{\"name\":\"ccccc\",\"description\":\"ccccc\",\"type\":2,\"trigger_time\":\"2021-08-20 05:54:54\",\"trigger_interval\":\"33\"}},{\"key\":4738053273,\"name\":\"新增数据\",\"description\":\"\",\"type\":6,\"props\":{\"columns\":[{\"16\":{\"type\":3,\"value\":\"11111\"}}],\"name\":\"1231321321\",\"description\":\"31221312\",\"target_id\":2}},{\"key\":3251668235,\"name\":\"结束\",\"description\":\"节点描述\",\"type\":9}]"
	b := []byte(str)

	var n Node

	err := json.Unmarshal(b, &n)
	t.Log(err)
}

func TestPoint(t *testing.T) {
	a := false
	p := &a
	tt(p)
	t.Log(a, p, *p)
}

func tt(bool2 *bool) {
	*bool2 = !(*bool2)
}
