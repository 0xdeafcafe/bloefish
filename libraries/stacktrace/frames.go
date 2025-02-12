package stacktrace

import (
	"fmt"
	"runtime"
	"strings"
)

type Frame struct {
	File string `json:"file"`
	Line int    `json:"line"`

	Function string `json:"function"`
}

// GetCallerFrames returns a slice of Frame objects representing the stack
//
// The skip argument is the number of frames to skip before starting to collect
// frames, with 0 identifying the frame for GetCallerFrames itself, and 1
// identifying the caller of GetCallerFrames
//
// No more than 100 frames will be collected
func GetCallerFrames(skip int) []Frame {
	stack := make([]uintptr, 100)
	length := runtime.Callers(skip+1, stack)
	stack = stack[:length]

	framesIter := runtime.CallersFrames(stack)
	frames := make([]Frame, 0, length)

	for {
		frame, more := framesIter.Next()

		frames = append(frames, Frame{
			File: frame.File,
			Line: frame.Line,

			Function: frame.Function,
		})

		if !more {
			break
		}
	}

	return frames
}

func (f Frame) String() string {
	return fmt.Sprintf("%s\n\t%s:%d", f.Function, f.File, f.Line)
}

func FormatFrames(frames []Frame) string {
	var sb strings.Builder

	for _, frame := range frames {
		sb.WriteString(frame.String())
		sb.WriteString("\n")
	}

	return sb.String()
}

// insert will put frame u into slice s at index at
func insert(s []Frame, u Frame, at int) []Frame {
	return append(s[:at], append([]Frame{u}, s[at:]...)...)
}

func MergeStacks(root []Frame, wrapped []Frame) []Frame {
	if len(wrapped) == 0 {
		return root
	}
	if len(wrapped) == 1 {
		return append(root, wrapped[0])
	}

	for idx, f := range root {
		if f == wrapped[0] {
			// root already contains the first frame of wrapped
			return root
		}
		if f == wrapped[1] {
			// Insert the first frame into the stack if the second frame is found
			return insert(root, wrapped[0], idx)
		}
	}
	return root
}
