package main

import "io"

func closeWith(closer io.Closer, handler func(error)) {
	err := closer.Close()
	if err != nil {
		handler(err)
	}
}
