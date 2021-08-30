package flow

type MAPS map[string]uint8

const (
	START_TABLE uint8 = iota + 1
	START_TIME
	DATA_ACT
	SERVICE
	BRANCH
	END
)

func (m MAPS) searchKey(key string) uint8 {
	for k, v := range m {
		if key == k {
			return v
		}
	}
	return 0
}
