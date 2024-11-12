package main

import (
	"network/modules/io"
)
func main(){
	newIO := io.Init[int]()
	newIO.Write([]int{1,2,3,4,5})
	io.PrintIndexedItem(newIO.Read())
	newIO.Write([]int{6,7,8,9,10})
	io.PrintIndexedItem(newIO.Read())
	newIO.Close()
	newIO = io.InitWithTimeOut[int](1)
	newIO.Write([]int{11,2,3,4,5})
	io.PrintIndexedItem(newIO.Read())
	newIO.Write([]int{6,7,8,9,10})
	io.PrintIndexedItem(newIO.Read())
	newIO.Close()

}