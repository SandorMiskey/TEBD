// region: packages

package log

import (
	"runtime"
)

// endregion: packages
// region: trace tells the details of the calling function

func Trace() (file string, line int, function string) {
	pc := make([]uintptr, 15)
	// n := runtime.Callers(2, pc)
	// frames := runtime.CallersFrames(pc[:n])
	// frame, _ := frames.Next()
	frame, _ := runtime.CallersFrames(pc[:runtime.Callers(2, pc)]).Next()
	file = frame.File
	line = frame.Line
	function = frame.Function
	return
}

// endregion: trace
