package utils

// Filter is a generic function that filters a slice of items based on a predicate function.
// It returns a new slice containing only the elements for which the predicate function returns true.
//
// Parameters:
//   - items: The slice of items to be filtered.
//   - pred: A function that takes an item of type T and returns a boolean value.
//
// Return:
//   - A new slice containing only the elements for which the predicate function returns true.
func Filter[T any](items []T, pred func(T) bool) []T {
	var res []T
	for _, v := range items {
		if pred(v) {
			res = append(res, v)
		}
	}
	return res
}

// Find is a generic function that finds the first item in a slice that satisfies a predicate function.
// It returns a pointer to the found item, or nil if no item satisfies the predicate.
//
// Parameters:
//   - items: The slice of items to be searched.
//   - pred: A function that takes an item of type T and returns a boolean value.
//
// Return:
//   - A pointer to the first item in the slice that satisfies the predicate, or nil if no item satisfies the predicate.
func Find[T any](items []T, pred func(T) bool) *T {
	for _, v := range items {
		if pred(v) {
			return &v
		}
	}
	return nil
}
