package utils

func ArrayHas[T comparable](arr []T, item T) bool {
	for _, v := range arr {
		if v == item {
			return true
		}
	}
	return false
}
