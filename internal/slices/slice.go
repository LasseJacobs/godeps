package slices

func Filter[T any](slice []T, filter func(T) bool) []T {
	var filtered []T
	for _, e := range slice {
		if filter(e) {
			filtered = append(filtered, e)
		}
	}
	return filtered
}
