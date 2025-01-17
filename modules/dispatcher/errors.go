package dispatcher

import (
	"errors"
	"fmt"

	"github.com/TreyBastian/colourize"
)

// custom errors for debugging

func backoffError(v bool, msg string) error {
	if !v {
		return errors.New(msg)
	}
	return nil
}

func argue(v bool, msg string) {
	if !v {
		error_msg := fmt.Sprintf(colourize.Colourize("[ERROR]: %v", colourize.Red), msg)
		panic(error_msg)
	}

}

func deferPanic(msg *chan interface{}) {
	if r := recover(); r != nil {
		(*msg) <- r
	}
}
