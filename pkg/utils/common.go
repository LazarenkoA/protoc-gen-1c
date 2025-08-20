package utils

func Ptr[T any](v T) *T {
	return &v
}

func Val[T any](v *T) T {
	var result T

	if v == nil {
		return result
	}

	return *v
}
