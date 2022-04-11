// region: packages

package log

import (
	"runtime"
)

// endregion: packages
// region: types

type Frame struct {
	File     string
	Function string
	Line     int
	More     bool
}

type Frames []Frame

// endregion: types
// region: trace tells the details of the calling function

func Trace(i ...int) Frames {

	// set depth and # of PCs
	depth := 2
	pcs := 10
	if len(i) > 0 {
		depth = i[0]
	}
	if len(i) > 1 {
		pcs = i[1]
	}

	// ask runtime.Callers for up to `pcs` PCs, excluding or including Trace() and runtime.Callers as a function of `depth`
	pc := make([]uintptr, pcs)
	n := runtime.Callers(depth, pc)
	if n == 0 {
		return nil // no PCs available, return now to avoid processing the zero Frame that would otherwise be returned by frames.Next below.
	}
	pc = pc[:n]
	frames := runtime.CallersFrames(pc)

	// loop to get frames, fixed number of PCs can expand to an indefinite number of Frames
	var stack Frames = make([]Frame, 0)
	for {
		frame, more := frames.Next()

		// process this frame
		stack = append(stack, Frame{
			File:     frame.File,
			Function: frame.Function,
			Line:     frame.Line,
			More:     more,
		})

		// check whether there are more frames to process after this one
		if !more {
			break
		}
	}
	return stack
}

// endregion: trace
