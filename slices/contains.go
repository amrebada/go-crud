package slices

type Primitives interface {
	~string | ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64 | ~bool
}

func Contains[T Primitives](slice []T, element T) bool {
	for _, e := range slice {
		if e == element {
			return true
		}
	}
	return false
}

func ContainsSlice[T Primitives](slice []T, elements []T) bool {
	for _, e := range elements {
		if !Contains(slice, e) {
			return false
		}
	}
	return true
}
