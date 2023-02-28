package util

// Concatenates two slices
func ConcatSlice[T any](first []T, second []T) []T {
	n := len(first)
	return append(first[:n:n], second...)
}
