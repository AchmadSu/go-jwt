package utils

import (
	"log"
	"runtime"
)

func logStackTrace(skip int) {
	const depth = 20
	pc := make([]uintptr, depth)
	n := runtime.Callers(skip, pc)
	frames := runtime.CallersFrames(pc[:n])

	log.Println("[STACK TRACE]")
	for {
		frame, more := frames.Next()
		log.Printf("  %s\n    %s:%d", frame.Function, frame.File, frame.Line)
		if !more {
			break
		}
	}
	log.Println("[END STACK TRACE]")
}
