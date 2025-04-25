package core

func Filter[T any](items []T, keep func(T) bool) []T {
	var out []T
	for _, x := range items {
		if keep(x) {
			out = append(out, x)
		}
	}
	return out
}

func Map[A any, B any](input []A, f func(A) B) []B {
	result := make([]B, 0, len(input))
	for _, a := range input {
		result = append(result, f(a))
	}
	return result
}
