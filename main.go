package main

import
//"encoding/binary"
//"fmt"
//"net"
//"os"
//"time"
//
//"github.com/brodo/go-opc/message"
//"github.com/brodo/go-opc/display"
//"github.com/brodo/go-opc/shader"

(
	"fmt"

	"./shader"
	//"github.com/brodo/go-opc/shader"
	//"github.com/brodo/go-opc/display"
)

func main() {
	fmt.Println("OPC client started")
	shader.Start()
	//
	/*conn, err := net.Dial("tcp", "10.23.42.141:7890")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not connect to server")
	}
	for {
		msg := message.EmptyMessage()
		//display.MakePixels()
		for y := 0; y < 32; y++ {
			for x := 0; x < 32; x++ {
				display.SetPixel(pixel.Pixel{1, 1, 1}, x, y)
				msg.SetData(display.GetBuffer().Bytes())
				binary.Write(conn, binary.BigEndian, msg.ToBytes())
				time.Sleep(100 * time.Millisecond)
			}
		}
		msg.SetData(display.GetBuffer().Bytes())
		binary.Write(conn, binary.BigEndian, msg.ToBytes())
		time.Sleep(4 * time.Millisecond)
	}*/
}
