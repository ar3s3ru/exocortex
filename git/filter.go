package git

func index(vs []string, t string) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}
	return -1
}

func include(vs []string, t string) bool {
	return index(vs, t) >= 0
}

func filter(vs []string, f func(string) bool) (vsf []string) {
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return
}
