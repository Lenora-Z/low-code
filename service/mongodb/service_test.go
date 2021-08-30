//Package mongodb Created by GoLand
//@User: lenora
//@Date: 2021/7/1
//@Time: 17:43

package mongodb

import (
	"testing"
)

func Test_Where(t *testing.T) {
	condition := []string{
		"aa", "bb", "cc",
	}
	key := "dfhjksd"

	mt := make(map[string]interface{})

	if len(condition) == 1 {
		mt[key] = condition[0]
	} else {
		mt[key] = In(condition)
	}
	t.Log(mt)
}
