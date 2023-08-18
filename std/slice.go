package std

func DelSliceElement[T comparable](s []T, e T) []T {
	dst := make([]T, len(s))
	for _, v := range s {
		if v != e {
			dst = append(dst, v)
		}
	}
	return dst
}

func DelSliceIndex[T comparable](s []T, idx int) []T {
	dst := make([]T, len(s))
	for i, v := range s {
		if i != idx {
			dst = append(dst, v)
		}
	}
	return dst
}
