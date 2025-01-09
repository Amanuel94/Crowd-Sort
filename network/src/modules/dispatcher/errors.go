package dispatcher

import (
	"errors"
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
		panic("DISPATCHER ERROR: " + msg)
	}

}

func deferPanic(msg *chan interface{}) {
	if r := recover(); r != nil {
		(*msg) <- r
	}
}
