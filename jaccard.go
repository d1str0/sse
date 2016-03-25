package sse

func Jaccard(set1, set2 []string) float64 {
	m := make(map[string]bool)
	for _, s := range set1 {
		m[s] = true
	}

	var intCount int
	for _, s2 := range set2 {
		if m[s2] {
			intCount++
		}
	}

	return float64(intCount) / float64(len(set1)+len(set2)-intCount)
}
