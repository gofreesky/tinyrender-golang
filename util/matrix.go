package util

import "fmt"

type Matrix struct {
	Col  int
	Row  int
	Data [][]float64
}

func NewEmpty(row int, col int) *Matrix {
	m := &Matrix{}
	m.Row = row
	m.Col = col
	m.Data = make([][]float64, row)
	for i := range m.Data {
		m.Data[i] = make([]float64, col)
	}
	return m
}

func NewIdentity(n int) *Matrix {
	m := NewEmpty(n, n)
	for i := 0; i < n; i++ {
		m.Data[i][i] = 1
	}
	return m
}

func NewFromSlice(f [][]float64) *Matrix {

	m := &Matrix{}
	if len(f) == 0 {
		return m
	}
	m.Row = len(f)
	m.Col = len(f[0])
	m.Data = f
	return m
}

func NewFromVector3(v Vector3) *Matrix {

	m := NewEmpty(4, 1)
	m.Data[0][0] = v.X
	m.Data[1][0] = v.Y
	m.Data[2][0] = v.Z
	m.Data[3][0] = 1
	return m
}

func (m *Matrix) ToVector3() Vector3 {
	if m.Data[3][0] == 0 {
		panic("m.Data[3][0] is zero")
	}
	w := m.Data[3][0]
	x := m.Data[0][0] / w
	y := m.Data[1][0] / w
	z := m.Data[2][0] / w
	return NewVec3(x, y, z)
}

func (m *Matrix) Print() {
	fmt.Printf("\n[")
	for i := 0; i < m.Row; i++ {
		for j := 0; j < m.Col; j++ {
			fmt.Printf("\t%.2f,", m.Data[i][j])
		}
		if i != m.Row-1 {
			fmt.Println("\t")
		}
	}
	fmt.Printf("]\n")
}

func (m *Matrix) Mul(to *Matrix) *Matrix {
	if m.Col != to.Row {
		panic("matrix row not equal col")
	}
	r := NewEmpty(m.Row, to.Col)
	for i := 0; i < r.Row; i++ {
		for j := 0; j < r.Col; j++ {
			for k := 0; k < m.Col; k++ {
				d1 := m.Data[i][k]
				d2 := to.Data[k][j]
				r.Data[i][j] += d1 * d2
			}
		}
	}
	return r
}

func (m *Matrix) Add(to *Matrix) *Matrix {
	if m.Row != to.Row || m.Col != to.Col {
		panic("matrix row/col not equal")
	}
	r := NewEmpty(m.Row, m.Col)
	for i := 0; i < r.Row; i++ {
		for j := 0; j < r.Col; j++ {
			r.Data[i][j] = m.Data[i][j] + to.Data[i][j]
		}
	}
	return r
}

func (m *Matrix) Sub(to *Matrix) *Matrix {
	if m.Row != to.Row || m.Col != to.Col {
		panic("matrix row/col not equal")
	}
	r := NewEmpty(m.Row, m.Col)
	for i := 0; i < r.Row; i++ {
		for j := 0; j < r.Col; j++ {
			r.Data[i][j] = m.Data[i][j] - to.Data[i][j]
		}
	}
	return r
}
