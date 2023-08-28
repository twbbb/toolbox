package maputil

// AddToMapMapStruct map的key是string,值也是map2
func AddToMapMapStruct[T any](m map[string]map[string]T, key1 string, key2 string, value T) {
	_, ok := m[key1]
	if !ok {
		m[key1] = map[string]T{}
	}
	m[key1][key2] = value
}
