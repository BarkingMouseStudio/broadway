package broadway

func concat(a, b []string) []string {
	c := make([]string, len(a)+len(b))
	copy(c, a)
	copy(c[len(a):], b)
	return c
}
