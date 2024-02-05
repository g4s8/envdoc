package main

import "fmt"

const debugLogs = false

func debug(f string, args ...any) {
	if !debugLogs {
		return
	}
	fmt.Printf("DEBUG: "+f+"\n", args...)
}
