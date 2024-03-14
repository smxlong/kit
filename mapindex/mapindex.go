package mapindex

// New builds a new map index for the given slice using the given key function.
//
// This converts a slice of the form:
//
//	[]T{T1, T2, ...}
//
// to a map of the form:
//
//	map[K]T{key(T1): T1, key(T2): T2, ...}
//
// The key function is used to extract the key from each element of the slice.
// The key must be a comparable type.
func New[T any, K comparable](slice []T, key func(T) K) map[K]T {
	m := make(map[K]T, len(slice))
	for _, v := range slice {
		m[key(v)] = v
	}
	return m
}

// Keys returns the keys of the map index. Do not depend on the order of the keys
// in the returned slice.
func Keys[K comparable, T any](m map[K]T) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Values returns the values of the map index. Do not depend on the order of the
// values in the returned slice.
func Values[K comparable, T any](m map[K]T) []T {
	values := make([]T, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

// Map returns a new slice containing the result of applying the function f to
// each element of the slice l.
func Map[T, U any](l []T, f func(T) U) []U {
	result := make([]U, len(l))
	for i, v := range l {
		result[i] = f(v)
	}
	return result
}

// Filter returns a new slice containing only the elements of the slice l for
// which the function f returns true.
func Filter[T any](l []T, f func(T) bool) []T {
	var result []T
	for _, v := range l {
		if f(v) {
			result = append(result, v)
		}
	}
	return result
}
