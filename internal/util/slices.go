package util

// Concatenates two slices
func ConcatSlice[T any](first []T, second []T) []T {
	n := len(first)
	return append(first[:n:n], second...)
}

func Contains[T string | int](s []T, e T) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
