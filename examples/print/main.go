package main

import (
	"github.com/acronis/journald"
)

func main() {
	journald.Print(journald.PriorityInfo, "Hello World!")
}
