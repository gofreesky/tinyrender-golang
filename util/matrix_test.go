package util

import "testing"

func Test001(t *testing.T) {
	f := [][]float64{
		{1, 2},
		{3, 4},
		{5, 6},
	}

	f2 := [][]float64{
		{1, 2, 3},
		{4, 5, 6},
	}

	m := NewFromSlice(f)
	m2 := NewFromSlice(f2)
	m3 := m.Mul(m2)
	m3.Print()
	m.Add(m).Print()
	m.Sub(m).Print()
}
