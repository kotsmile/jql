package util

func At[T any](index int, slice []T) (T, bool) {
	var t T
	if index < len(slice) {
		t = slice[index]
		return t, true
	}
	return t, false
}

func Next[T any](slice *[]T) (T, bool) {
	el, ok := Peek(*slice)
	if ok {
		*slice = (*slice)[1:]
	}
	return el, ok
}

func Peek[T any](slice []T) (T, bool) {
	var t T
	if len(slice) == 0 {
		return t, false
	}

	t = slice[0]
	return t, true
}
