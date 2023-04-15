package stacktrace

import (
	"runtime"
)

type (
	RawTrace      = uintptr
	RawStackTrace = []RawTrace

	Trace struct {
		File string `json:"file"`
		Func string `json:"func"`
		Line uint   `json:"line"`
	}

	StackTrace = []Trace
)

func Fix(indent uint) RawStackTrace {
	stack := make(RawStackTrace, 32)
	n := runtime.Callers(int(indent), stack)
	return stack[:n]
}

func Parse(rawStackTrace RawStackTrace) StackTrace {
	return parseRawStackTrace(rawStackTrace)
}

func parseRawStackTrace(rawStackTrace RawStackTrace) StackTrace {
	l := len(rawStackTrace)
	if l == 0 {
		return nil
	}

	parsed := make(StackTrace, 0, l)
	for _, rawTrace := range rawStackTrace {
		parsed = append(parsed, parseRawTrace(rawTrace))
	}

	return parsed
}

func parseRawTrace(rawTrace RawTrace) Trace {
	fn := runtime.FuncForPC(rawTrace)
	if fn == nil {
		return Trace{
			File: "unknown",
		}
	}

	file, line := fn.FileLine(rawTrace)
	return Trace{
		Func: fn.Name(),
		File: file,
		Line: uint(line),
	}
}
