package selector

// custom errors for debugging

func argue(v bool, msg string) {
	if !v {
		panic(msg)
	}
}

func deferPanic(msg *chan interface{}) {
	if r := recover(); r != nil {
		(*msg) <- r
	}
}
