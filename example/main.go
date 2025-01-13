package main

import (
	"fmt"

	"github.com/amanuel94/crowdsort/network/modules/io"
)

func main() {

	io := io.Init[int64]()
	fmt.Println(io)

}
