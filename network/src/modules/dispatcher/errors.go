package dispatcher

import (
	"errors"
)

// custom errors for debugging

func argue(v bool, msg string) error {
	if !v {
		return errors.New(msg)
	}
	return nil
}
