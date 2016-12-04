package pixel

import (
	"bytes"
	"encoding/binary"
)

type Pixel struct {
	Red   float64
	Green float64
	Blue  float64
}

func (p Pixel) Multiply(scalar float64) Pixel {
	return Pixel{p.Red * scalar, p.Green * scalar, p.Blue * scalar}
}

func (p1 Pixel) Add(p2 Pixel) Pixel {
	return Pixel{p1.Red + p2.Red, p1.Green + p2.Green, p1.Blue + p2.Blue}
}

func Interpolate(left, right Pixel, t float64) Pixel {
	return left.Multiply(1 - t).Add(right.Multiply(t))
}

func (p Pixel) ToBuffer() *bytes.Buffer {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, byte(p.Red * 255))
	binary.Write(buf, binary.BigEndian, byte(p.Green * 255))
	binary.Write(buf, binary.BigEndian, byte(p.Blue * 255))
	return buf
}
