package dbg

import "fmt"

// Log simply logs any number of objects. This method is simply a fire and forget wrapper for fmt.Println
func Log(logItems ...interface{}) {
	fmt.Print("[.] ")
	fmt.Println(logItems...)
}

// DLog stands for Delimited Log and logs any number of objects with a delimiter between them.
// This is very useful if you want to print something like "VarA: 32 | VarB: 33"
func DLog(delim string, logItems ...interface{}) {
	paramCount := len(logItems)

	fmt.Print("[.] ")

	for ind, obj := range logItems {
		if ind == paramCount-1 {
			fmt.Println(obj)
		} else {
			fmt.Print(obj)
			fmt.Print(delim)
		}
	}
}

// LogError prints out a distinctly formatted error message so the program may continue to run, but the error is made known to the user.
func LogError(logItems ...interface{}) {
	fmt.Print("[X] ")
	fmt.Println(logItems...)
}
