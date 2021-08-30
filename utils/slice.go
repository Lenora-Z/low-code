//Created by Goland
//@User: lenora
//@Date: 2021/3/9
//@Time: 9:51 下午
package utils

func IntersectionInt(oldIds, newIds []uint64) []uint64 {
	var intersection []uint64
	for _, oldId := range oldIds {
		for _, newId := range newIds {
			if oldId == newId {
				intersection = append(intersection, oldId)
				break
			}
		}
	}
	return intersection
}

func IsContain(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

func DifferenceUInt64(oldIds, newIds []uint64) []uint64 {
	difference := make([]uint64, 0, len(oldIds))
	for _, oldId := range oldIds {
		if isContain := IsContainUInt64(newIds, oldId); isContain == false {
			difference = append(difference, oldId)
		}
	}
	return difference
}


func IsContainUInt64(items []uint64, item uint64) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}