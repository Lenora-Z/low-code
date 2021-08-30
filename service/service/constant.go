//Created by Goland
//@User: lenora
//@Date: 2021/3/11
//@Time: 2:59 下午
package service

const (
	PENDING uint8 = iota
	VALID
	INVALID
)

const (
	IN uint8 = iota + 1
	OUT
)

const (
	CHAIN_UP int8 = iota + 1
	EMAIL
	PACKAGE
)
