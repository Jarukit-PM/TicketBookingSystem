package httputil

// JSONSlice ensures nil slices encode as [] instead of null in JSON responses.
func JSONSlice[T any](s []T) []T {
	if s == nil {
		return []T{}
	}
	return s
}
