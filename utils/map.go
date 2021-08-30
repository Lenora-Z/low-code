package utils

func GetStringKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for i := range m {
		keys = append(keys, i)
	}
	return keys
}
