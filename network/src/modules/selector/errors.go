package selector

// custom errors for debugging


func argue(v bool, msg string){
	if !v{
		panic(msg)
	}
}