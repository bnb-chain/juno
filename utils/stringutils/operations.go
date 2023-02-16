package stringutils

// Difference finds the Set-wise Difference: {one} - {others}
func Difference(one []string, others []string) (diff []string) {
	m := make(map[string]bool)

	for _, item := range others {
		m[item] = true
	}

	for _, item := range one {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}
	return
}
