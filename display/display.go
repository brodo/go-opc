package display

import (
	"github.com/brodo/go-opc/pixel"
	"bytes"
	"math"
)


var (
	tick int
	pixels = make([]pixel.Pixel, 512)
	t float64
	left = pixel.Pixel{0,0,0}
 	right = pixel.Pixel{1,1,1}
)

func MakePixels() {
	interval := 1.0 / 512.0
	for x := 0; x < 32; x ++ {
		for y := 0; y < 32; y ++ {
			t = t + interval
			SetPixel(pixel.Interpolate(left, right, t), x, y)
		}
	}
}

// x = pos % 16
// y = pos % 16
// y * 16 + x = arraypos

func SetPixel(p pixel.Pixel, x, y int){
	if int(math.Pow(-1, float64(x + y))) == -1 {
		mappedX := (x - (1 - y % 2)) / 2
		//mappedY := (y - (1 - x % 2)) / 2
		pixels[y * 16 + mappedX] = p
	}
}

func GetBuffer() *bytes.Buffer {
	buf := new(bytes.Buffer)
	for _, p := range(pixels){
		buf.Write(p.ToBuffer().Bytes())
	}
	return buf
}
