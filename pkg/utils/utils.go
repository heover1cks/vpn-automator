package utils

import "sort"

type TypesForValue interface {
	string | int | int64 | float32 | float64
}

// TODO: apply go generic
//func SortedKeys[K comparable, V TypesForValue](m map[K]V) []K {
//	keys := make([]K, len(m))
//	i := 0
//	for k := range m {
//		keys[i] = k
//		i++
//	}
//	sort.Strings(keys)
//	return keys
//}

func SortedStringKeys(m map[string]string) []string {
	keys := make([]string, len(m))
	idx := 0
	for k := range m {
		keys[idx] = k
		idx++
	}
	sort.Strings(keys)
	return keys
}
