//Created by Goland
//@User: lenora
//@Date: 2021/1/23
//@Time: 11:06 上午
package utils

import "strconv"

type StrStruct struct {
	str int
}

func NewStr(str string) *StrStruct {
	r := new(StrStruct)
	r.str, _ = strconv.Atoi(str)
	return r
}

func (obj *StrStruct) Int() int {
	return obj.str
}

func (obj *StrStruct) Int8() int8 {
	return int8(obj.str)
}

func (obj *StrStruct) Int16() int16 {
	return int16(obj.str)
}

func (obj *StrStruct) Int32() int32 {
	return int32(obj.str)
}

func (obj *StrStruct) Int64() int64 {
	return int64(obj.str)
}

func (obj *StrStruct) Uint() uint {
	return uint(obj.str)
}

func (obj *StrStruct) Uint8() uint8 {
	return uint8(obj.str)
}

func (obj *StrStruct) Uint16() uint16 {
	return uint16(obj.str)
}

func (obj *StrStruct) Uint32() uint32 {
	return uint32(obj.str)
}

func (obj *StrStruct) Uint64() uint64 {
	return uint64(obj.str)
}

func Strlen(str string) int {
	return len([]rune(str))
}
