package std

func DelSliceIndex[T any](src []T, idx int) []T {
	dst := make([]T, 0, len(src)-1)
	for i, v := range src {
		if i != idx {
			dst = append(dst, v)
		}
	}
	return dst
}

func DelSliceElement[T comparable](src []T, e T) []T {
	dst := make([]T, 0, len(src)-1)
	for _, v := range src {
		if v != e {
			dst = append(dst, v)
		}
	}
	return dst
}
