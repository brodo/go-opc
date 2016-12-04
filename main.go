package main

import (
	//"encoding/binary"
	//"fmt"
	//"net"
	//"os"
	//"time"
	//
	//"github.com/brodo/go-opc/message"
	//"github.com/brodo/go-opc/display"
	//"github.com/brodo/go-opc/shader"



	"fmt"
	"github.com/brodo/go-opc/shader"
)


func main() {
	fmt.Println("OPC client started")
	shader.Start()
	//
	//conn, err := net.Dial("tcp", "10.23.42.141:7890")
	//if err != nil {
	//	fmt.Fprintln(os.Stderr, "Could not connect to server")
	//}
	//for {
	//	msg := message.EmptyMessage()
	//	display.MakePixels()
	//	msg.SetData(display.GetBuffer().Bytes())
	//	binary.Write(conn, binary.BigEndian, msg.ToBytes())
	//	time.Sleep(4 * time.Millisecond)
	//}
}
