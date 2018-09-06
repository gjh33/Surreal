package math

// Matrix interface represents a standard matrix implementation in the Surreal Engine
type Matrix interface {
	Get(row int, col int) float32        // Returns the element at the corresponding row and column. Indexes start at 0 you mongrel. Values too large or too small should wrap.
	Set(row int, col int, value float32) // Sets the element at the corresponding row and column. Indexes start at 0 you mongrel. Values too large or too small should wrap.
	// TODO: GetSubMatrix(x1 int, y1 int, x2 int, y2 int) Matrix
	// TODO: SetSubMatrix(x int, y int, sub *Matrix)
	NumRows() int             // Returns the number of rows in the matrix
	NumCols() int             // Returns the number of columns in the matrix
	ColMajorData() *[]float32 // Returns the data in column major order for use with open gl
}
