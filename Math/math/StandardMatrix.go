package math

// StandardMatrix implements the Matrix interface and represents a uncompressed, unoptimized matrix.
// Stored in column major order
// TODO: Deprecate this garbage
type StandardMatrix struct {
	data    []float32
	numRows int
	numCols int
}

// CreateStandardMatrix is the standard constructor for a Matrix4x4
func CreateStandardMatrix(data []float32, numRows int, numCols int) *StandardMatrix {
	mat := new(StandardMatrix)
	mat.data = data[:numCols*numRows]
	mat.numRows = numRows
	mat.numCols = numCols
	return mat
}

// StandardMatrixZeros returns a matrix of zeros. It's an alternate constructor
func StandardMatrixZeros(numRows int, numCols int) *StandardMatrix {
	mat := new(StandardMatrix)
	mat.data = make([]float32, numRows*numCols, numRows*numCols)
	mat.numRows = numRows
	mat.numCols = numCols
	return mat
}

// StandardMatrixIdentity returns an identity matrix of arbitrary compisition
func StandardMatrixIdentity(numRows int, numCols int) *StandardMatrix {
	mat := StandardMatrixZeros(numRows, numCols)
	for i, j := 0, 0; i < numRows && j < numCols; i, j = i+1, j+1 {
		mat.Set(i, j, 1)
	}
	return mat
}

// Get is Implementing Matrix Interface
func (mat *StandardMatrix) Get(row int, col int) float32 {
	row = wrapIndex(row, mat.NumRows())
	col = wrapIndex(col, mat.NumCols())
	index := mat.NumRows()*col + row
	return mat.data[index]
}

// Set is Implementing Matrix Interface
func (mat *StandardMatrix) Set(row int, col int, value float32) {
	row = wrapIndex(row, mat.NumRows())
	col = wrapIndex(col, mat.NumCols())
	index := mat.NumRows()*col + row
	mat.data[index] = value
}

// NumRows is Implementing Matrix Interface
func (mat *StandardMatrix) NumRows() int {
	return 4
}

// NumCols is Implementing Matrix Interface
func (mat *StandardMatrix) NumCols() int {
	return 4
}

// MulM multiplies this matrix by another via this * other and returns a pointer to a NEW matrix as the result
// TODO: Optimize this garbage
// TODO: Maybe don't return a new matrix because composite operations start to consume fuck loads of memory
// TODO: Can this be moved to interface and maybe return a Matrix instead of *StandardMatrix?
func (mat *StandardMatrix) MulM(other Matrix) *StandardMatrix {
	if mat.NumCols() != other.NumRows() || mat.NumRows() != other.NumCols() {
		panic("Invalid Operation: Matrix dimensions are incompatible for inner product")
	}
	toRet := StandardMatrixZeros(mat.NumRows(), other.NumCols())
	for i1 := 0; i1 < mat.NumRows(); i1++ {
		for j2 := 0; j2 < other.NumCols(); j2++ {
			var sum float32
			for ind := 0; ind < mat.NumCols(); ind++ {
				sum += mat.Get(i1, ind) * (other.Get(ind, j2))
			}
			toRet.Set(i1, j2, sum)
		}
	}

	return toRet
}

// ColMajorData implements the Matrix interface
func (mat *StandardMatrix) ColMajorData() *[]float32 {
	return &mat.data
}

// Returns the correctly wrapped index for a list with a count number of items
func wrapIndex(rawIndex int, count int) int {
	return (count + (rawIndex % count)) % count
}
