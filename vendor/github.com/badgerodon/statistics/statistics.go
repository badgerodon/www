package statistics

import (
	"math"

	. "github.com/badgerodon/lalg"
)

// Calculate the covariance between two return series
func Covariance(x, y Vector) float64 {
	if len(x) != len(y) {
		panic("Vector lengths must be the same")
	}

	n := len(x)
	sum, xsum, xmean, ysum, ymean := 0.0, 0.0, 0.0, 0.0, 0.0

	for i := 0; i < n; i++ {
		xsum += x[i]
		ysum += y[i]
	}
	xmean = xsum / float64(n)
	ymean = ysum / float64(n)

	for i := 0; i < n; i++ {
		sum += (x[i] - xmean) * (y[i] - ymean)
	}

	return sum / float64(n-1)
}

// Calculate the covariance matrix of a matrix of returns
// Each return series should be placed in the row:
//  series1: 1, 2, 3, 4, 5
//  series2: 1, 2, 3, 4, 5
func CovarianceMatrix(mat Matrix) Matrix {
	n := mat.Rows
	c := NewMatrix(n, n)

	for i := 0; i < n; i++ {
		// variance on the diagonal
		c.Set(i, i, Variance(mat.GetRow(i)))
		for j := 0; j < i; j++ {
			cell := Covariance(mat.GetRow(i), mat.GetRow(j))
			// Same across the diagonal
			c.Set(i, j, cell)
			c.Set(j, i, cell)
		}
	}

	return c
}

// Make a matrix positive definite by removing rows which are too similar
func MakePositiveDefinite(mat Matrix) ([]int, Matrix) {
	rowIndices := make([]int, mat.Rows)
	mod := NewMatrix(mat.Rows, mat.Cols)

	for i := 0; i < mat.Rows; i++ {
		rowIndices[i] = i
		for j := 0; j < mat.Cols; j++ {
			mod.Set(i, j, mat.Get(i, j))
		}
	}

	ε := 0.00000001
	last := 0
	for j := 0; j < mat.Cols; j++ {
		i := last
		for i < mat.Rows && math.Abs(mod.Get(j, i)) < ε {
			i++
		}

		if i < mat.Rows {
			if math.Abs(mod.Get(i, j)) >= ε {
				// Swap the indices
				rowIndices[i], rowIndices[last] = rowIndices[last], rowIndices[i]

				// Swap the rows
				for k := 0; k < mat.Cols; k++ {
					t := mod.Get(i, k)
					mod.Set(i, k, mod.Get(last, k))
					mod.Set(last, k, t)
				}

				// Multiply out the remaining ros so that they have zeros in j
				for i = last + 1; i < mat.Rows; i++ {
					factor := mod.Get(i, j)
					for k := 0; k < mat.Cols; k++ {
						v := (factor / mod.Get(last, j)) * mod.Get(last, k)
						mod.Set(i, k, mod.Get(i, k)-v)
					}
				}

				last++
			}
		}
	}

	indices := make([]int, last)
	final := NewMatrix(last, last)
	for i := 0; i < last; i++ {
		indices[i] = rowIndices[i]
		for j := 0; j < last; j++ {
			final.Set(i, j, mat.Get(rowIndices[i], rowIndices[j]))
		}
	}
	return indices, final
}

// Take an absolute return series and transform it into a relative one
// The new vector will have one less item
func Relativize(vector Vector) Vector {
	nv := NewVector(len(vector) - 1)
	for i := 1; i < len(vector); i++ {
		nv[i-1] = (vector[i] - vector[i-1]) / vector[i-1]
	}
	return nv
}

// Find the variance of a vector
func Variance(vector Vector) float64 {
	n := 0.0
	mean := 0.0
	S := 0.0
	delta := 0.0

	for _, v := range vector {
		n++
		delta = v - mean
		mean = mean + (delta / n)
		S += delta * (v - mean)
	}

	return S / (n - 1)
}
