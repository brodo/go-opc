package message

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// A OPC message
type Message struct {
	channel byte
	command byte
	length  uint16
	data    []byte
}

// Convert a message to bytes which can be sent to an OPC server
func (msg Message) ToBytes() []byte {
	result := make([]byte, msg.length+4)
	result[0] = msg.channel
	result[1] = msg.command
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, msg.length)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	result[2], err = buf.ReadByte()
	result[3], err = buf.ReadByte()
	if err != nil {
		fmt.Println("binary.ReadByte failed:", err)
	}
	for i := 0; i < int(msg.length); i++ {
		result[i+4] = msg.data[i]
	}
	return result
}

func (msg *Message) SetData(data []byte) {
	msg.data = data
}

func EmptyMessage() Message {
	data := make([]byte, 512 * 3)
	return Message{0, 0, 512 * 3, data}
}

