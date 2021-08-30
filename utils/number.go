//Created by Goland
//@User: lenora
//@Date: 2021/3/12
//@Time: 2:21 下午
package utils

import (
	"fmt"
	"math/rand"
	"strings"
)

func RandInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func UintJoin(group []uint64, sep string) string {
	list := make([]string, 0, len(group))
	for _, x := range group {
		list = append(list, fmt.Sprintf("%d", x))
	}
	return strings.Join(list, sep)
}
