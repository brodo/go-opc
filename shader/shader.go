package shader

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"image"
	"image/draw"
	_ "image/png"
	"log"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/brodo/go-opc/display"
	"github.com/brodo/go-opc/message"
	"github.com/brodo/go-opc/pixel"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Texture struct {
	TexId         uint32
	TextureNumber uint32
	ChannelNumber uint32
}

const SCALE_X = 4
const SCALE_Y = 4

const RES_X = 32
const RES_Y = 32

const FB_RES_X = SCALE_X * RES_X
const FB_RES_Y = SCALE_Y * RES_Y

var pixels = make([]byte, FB_RES_X*FB_RES_Y*4)
var textureCount = uint32(0)
var textures = make([]Texture, 0)

func init() {
	runtime.LockOSThread()
}

func Start(shaderFile string) {
	// Initialize glfw
	if err := glfw.Init(); err != nil {
		log.Fatal(err)
	}
	defer glfw.Terminate()

	// Window hints
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	// Create window
	window, err := glfw.CreateWindow(FB_RES_X, FB_RES_Y, "Window", nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	window.MakeContextCurrent()

	// Initialize gl
	err = gl.Init()
	if err != nil {
		log.Fatal(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	// Triangle vertices
	vertexBufferData := []float32{
		-1.0, 1.0, 0.0,
		1.0, 1.0, 0.0,
		-1.0, -1.0, 0.0,
		-1.0, -1.0, 0.0,
		1.0, -1.0, 0.0,
		1.0, 1.0, 0.0}

	// Create Vertex array object
	var vertexArrayID uint32
	gl.GenVertexArrays(1, &vertexArrayID)
	gl.BindVertexArray(vertexArrayID)

	var vertexBuffer uint32
	gl.GenBuffers(1, &vertexBuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertexBufferData)*4, gl.Ptr(vertexBufferData), gl.STATIC_DRAW)

	gl.ClearColor(1.0, 1.0, 1.0, 1.0)
	// load shaders
	prog, err := NewProgram("shader/vertex.glsl", shaderFile)
	if err != nil {
		log.Fatal(err)
	}

	files, _ := filepath.Glob("*")
	for _, file := range files {
		r, _ := regexp.Compile("\\.chan([0-9]+)")
		matches := r.FindStringSubmatch(file)

		if len(matches) > 0 {
			fmt.Println(matches[1])
			if number, err := strconv.Atoi(matches[1]); err == nil {
				textures = append(textures, MakeTexture(file, uint32(number)))
			}
		}
	}

	conn, err := net.Dial("tcp", "10.23.42.141:7890")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not connect to server")
	}

	gl.ClearColor(0.11, 0.545, 0.765, 0.0)

	t0 := NanosNow()

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(prog)

		iResolutionIndex := gl.GetUniformLocation(prog, gl.Str("iResolution\x00"))
		if iResolutionIndex >= 0 {
			gl.Uniform2f(iResolutionIndex, float32(FB_RES_X), float32(FB_RES_Y))
		}

		iGlobalTimeIndex := gl.GetUniformLocation(prog, gl.Str("iGlobalTime\x00"))
		if iGlobalTimeIndex >= 0 {
			gl.Uniform1f(iGlobalTimeIndex, float32(NanosNow()-t0)/1000000000.0)
		}

		for _, texture := range textures {
			iChannelNumber := texture.ChannelNumber
			iChannelStr := fmt.Sprintf("iChannel%d\x00", iChannelNumber)
			iChannelIndex := gl.GetUniformLocation(prog, gl.Str(iChannelStr))

			if iChannelIndex >= 0 {
				gl.ActiveTexture(gl.TEXTURE0 + texture.TextureNumber)
				gl.BindTexture(gl.TEXTURE_2D, texture.TexId)
				gl.Uniform1i(iChannelIndex, int32(texture.TextureNumber))
			}
		}

		gl.EnableVertexAttribArray(0)

		gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

		gl.DrawArrays(gl.TRIANGLES, 0, 6)

		gl.DisableVertexAttribArray(0)

		gl.ReadPixels(0, 0, FB_RES_X, FB_RES_Y, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(pixels))

		for x := 0; x < RES_X; x++ {
			for y := 0; y < RES_Y; y++ {
				buf_pos := y*SCALE_Y*FB_RES_X*4 + x*SCALE_X*4

				pixel := pixel.Pixel{
					Red:   float64(pixels[buf_pos]) / 255.0,
					Green: float64(pixels[buf_pos+1]) / 255.0,
					Blue:  float64(pixels[buf_pos+2]) / 255.0,
				}

				display.SetPixel(pixel, x, y)
			}
		}

		msg := message.EmptyMessage()
		msg.SetData(display.GetBuffer().Bytes())
		binary.Write(conn, binary.BigEndian, msg.ToBytes())

		time.Sleep(4 * time.Millisecond)

		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func NanosNow() int64 {
	return time.Now().UnixNano()
}

// Create a new program to run. Requires path to vertex shader and fragment
// shader files
func NewProgram(vertexFilePath, fragmentFilePath string) (uint32, error) {
	// Load both shaders
	vertexShaderID, fragmentShaderID, err := LoadShaders(vertexFilePath, fragmentFilePath)
	if err != nil {
		return 0, err
	}

	// Create new program
	programID := gl.CreateProgram()
	gl.AttachShader(programID, vertexShaderID)
	gl.AttachShader(programID, fragmentShaderID)
	gl.LinkProgram(programID)

	// Check status of program
	var status int32
	gl.GetProgramiv(programID, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(programID, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(programID, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	// Detach shaders
	gl.DetachShader(programID, vertexShaderID)
	gl.DetachShader(programID, fragmentShaderID)

	// Delete shaders
	gl.DeleteShader(vertexShaderID)
	gl.DeleteShader(fragmentShaderID)

	return programID, nil
}

// Load both shaders and return
func LoadShaders(vertexFilePath, fragmentFilePath string) (uint32, uint32, error) {

	// Compile vertex shader
	vertexShaderID, err := CompileShader(ReadShaderCode(vertexFilePath), gl.VERTEX_SHADER)
	if err != nil {
		return 0, 0, nil
	}

	// Compile fragment shader
	fragmentShaderID, err := CompileShader(ReadShaderCode(fragmentFilePath), gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, 0, nil
	}

	return vertexShaderID, fragmentShaderID, nil
}

// Compile shader. Source is null terminated c string. shader type is self
// explanatory
func CompileShader(source string, shaderType uint32) (uint32, error) {
	// Create new shader
	shader := gl.CreateShader(shaderType)
	// Convert shader string to null terminated c string
	shaderCode, free := gl.Strs(source)
	defer free()

	gl.ShaderSource(shader, 1, shaderCode, nil)

	// Compile shader
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		fmt.Errorf("Failed to compile %v: %v", source, log)
	}

	return shader, nil
}

// Read shader code from file
func ReadShaderCode(filePath string) string {
	code := ""

	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		code += "\n" + scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	code += "\x00"

	return code
}

func MakeTexture(file string, channelNumber uint32) Texture {
	imgFile, err := os.Open(file)
	if err != nil {
		log.Fatalf("texture %q not found on disk: %v\n", file, err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		panic(err)
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		panic("Unsupported stride!")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var texId uint32
	gl.Enable(gl.TEXTURE_2D)
	gl.GenTextures(1, &texId)
	gl.BindTexture(gl.TEXTURE_2D, texId)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	texture := Texture{
		TexId:         texId,
		TextureNumber: textureCount,
		ChannelNumber: channelNumber,
	}

	textureCount++

	return texture
}
