package main

import (
	"network/modules/io"
)
func main(){
	newIO := io.Init[int64]()
	newIO.WriteInt([]int64{1,2,3,4,5})
	io.PrintIndexedItem(newIO.Read())
	newIO.WriteInt([]int64{6,7,8,9,10})
	io.PrintIndexedItem(newIO.Read())
	newIO.Close()
	newIO = io.InitWithTimeOut[int64](1)
	newIO.WriteInt([]int64{11,2,3,4,5})
	io.PrintIndexedItem(newIO.Read())
	newIO.WriteInt([]int64{6,7,8,9,10})
	io.PrintIndexedItem(newIO.Read())
	newIO.Close()

}