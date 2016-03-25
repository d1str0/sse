package sse

import "testing"

func TestJaccard(t *testing.T) {
	a := []string{"a", "b", "c"}
	b := []string{"a", "b", "c"}

	c := Jaccard(a, b)
	t.Logf("a, b: %f\n", c)
}
