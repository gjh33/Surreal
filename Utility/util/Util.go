package util

import (
	"path/filepath"

	"github.com/kardianos/osext"
)

// DataFolderName is the folder from surreal root for accessing data
const DataFolderName string = "ext"

// ExecutablesFolderName is the folder from surreal root for all executables
const ExecutablesFolderName string = "bin"

var surrealRoot string

func init() {
	execRoot, err := osext.ExecutableFolder()
	if err != nil {
		panic("Failed to find Surreal Root Directory")
	}

	sr, err := filepath.Abs(execRoot + "/..")

	surrealRoot = sr
}

// SurrealRoot returns a copy of the absolute path to the root of the Surreal installation directory.
// WARNING: There is currently no guarenteed solution to this in Go. SymLinks may fail.
func SurrealRoot() string {
	return string([]byte(surrealRoot))
}

// DataRoot returns a copy of the absolute path to the folder for external program data
func DataRoot() string {
	return filepath.Join(surrealRoot, DataFolderName, "Surreal")
}

// ExecutablesRoot returns a copy of the absolute path to the folder for all executables
func ExecutablesRoot() string {
	return filepath.Join(surrealRoot, ExecutablesFolderName)
}
