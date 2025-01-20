package selector

import (
	"fmt"

	"github.com/TreyBastian/colourize"
)

// custom errors for debugging

func argue(v bool, msg string) {
	if !v {
		panic(msg)
	}
}

// TODO: Make this better
func deferPanic(msg *chan interface{}) {
	if r := recover(); r != nil {
		error_msg := fmt.Sprintf(colourize.Colourize("[ERROR]: %v", colourize.Red), r)
		(*msg) <- error_msg
	}
}
