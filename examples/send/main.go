package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/acronis/journald"
)

func main() {
	err := journald.Send("Hello World!", journald.PriorityInfo, map[string]interface{}{
		"HOME":        os.Getenv("HOME"),
		"TERM":        os.Getenv("TERM"),
		"N_GOROUTINE": runtime.NumGoroutine(),
		"N_CPUS":      runtime.NumCPU(),
		"TRACE":       runtime.ReadTrace(),
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
