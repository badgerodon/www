package lalg 

import (
	"fmt"
)

type (
	Matrix struct {
		Elements []float64
		Rows, Cols int
	}
	
	Vector []float64
)
/**
 * Create a new matrix
 */
func NewMatrix(rows, cols int) Matrix {
	return Matrix{make([]float64, rows*cols),rows,cols}
}
/**
 * Create a new identity matrix of size "size"
 */
func NewIdentity(size int) Matrix {
	m := NewMatrix(size, size)
	for i := 0; i < size; i++ {
		m.Set(i, i, 1.0)
	}
	return m
}

// Get an element from the matrix
func (this Matrix) Get(row, col int) float64 {
	return this.Elements[row * this.Cols + col]
}
// Get a column from the matrix
func (this Matrix) GetCol(col int) []float64 {
	c := make([]float64, this.Rows)
	for i := 0; i < this.Rows; i++ {
		c[i] = this.Elements[i * this.Cols + col]
	}
	return c
}
// Get a row from the matrix
func (this Matrix) GetRow(row int) []float64 {
	return this.Elements[(row * this.Cols) : ((row+1) * this.Cols)]
}
// Set an element in the matrix
func (this Matrix) Set(row, col int, value float64) {
	this.Elements[row * this.Cols + col] = value
}
// Pretty-print the matrix
func (this Matrix) String() string {
	str := ""
	for i := 0; i < this.Rows; i++ {
		if i > 0 {		
			str += "\n"
		}
		for j := 0; j < this.Cols; j++ {
			str += fmt.Sprintf("%6.2f", this.Get(i,j))
		}
	}
	return str
}
// create a new vector
func NewVector(size int) Vector {
	return make([]float64, size)
}