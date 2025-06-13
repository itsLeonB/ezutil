package ezutil

// MapSlice applies a mapping function to each element of an input slice and returns a new slice.
// The function transforms elements of type T to type U using the provided mapperFunc.
// This is a generic utility for functional-style slice transformations.
func MapSlice[T any, U any](input []T, mapperFunc func(T) U) []U {
	output := make([]U, len(input))

	for i, v := range input {
		output[i] = mapperFunc(v)
	}

	return output
}