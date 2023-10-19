package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
	"os"
	"strings"
	"sync"

	vidio "github.com/AlexEidt/Vidio"
)

//bytesperframe: sqrt(bytesperframe*8) has to be a whole number
var bytesperframe int = 512

var framesize int = 1024

var bitsPerFrame int = bytesperframe * 8
var bitsPerRow int = int(math.Sqrt(float64(bitsPerFrame)))

var bitPixelSize int = int(math.Floor(float64(framesize / bitsPerRow)))

func encodeDataIntoPicture(data []byte) *image.RGBA {
	resultimg := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{framesize, framesize}})

	n_row := 0
	n_col := 0
	for byte_i := 0; byte_i < len(data)/3; byte_i++ {
		for i := 7; i >= 0; i-- {
			bit_r := (data[byte_i] >> i) & 1
			bit_g := (data[byte_i+bytesperframe] >> i) & 1
			bit_b := (data[byte_i+(bytesperframe*2)] >> i) & 1
			if n_row > bitsPerRow-1 {
				n_row = 0
				n_col++
			}
			rect := image.Rect(n_row*bitPixelSize, n_col*bitPixelSize, n_row*bitPixelSize+bitPixelSize, n_col*bitPixelSize+bitPixelSize)
			r := uint8(255 * bit_r)
			g := uint8(255 * bit_g)
			b := uint8(255 * bit_b)
			draw.Draw(resultimg, rect.Bounds(), &image.Uniform{color.RGBA{r, g, b, 255}}, image.ZP, draw.Src)
			n_row++
		}
	}
	return resultimg
}

var mutexwrite = &sync.Mutex{}

func readFileWriteFrames(inputfilename string, videofilename string) {
	dat, _ := os.Open(inputfilename)
	stat, _ := dat.Stat()
	defer dat.Close()
	var wg sync.WaitGroup
	var images map[int]*image.RGBA = make(map[int]*image.RGBA)
	options := vidio.Options{
		FPS: 25,
	}
	writer, _ := vidio.NewVideoWriter(videofilename, framesize, framesize, &options)

	tempvideoi := 0

	for i := 0; i < int(math.Ceil(float64(stat.Size())/float64(bytesperframe)/3)); i++ {
		wg.Add(1)
		go writeDataToFrame(dat, i, &wg, images, tempvideoi)
		tempvideoi++
		if i%100 == 0 {
			fmt.Println(i)
			wg.Wait()
			for v_i := 0; v_i < tempvideoi; v_i++ {
				writer.Write(images[v_i].Pix)
			}
			tempvideoi = 0
		}
	}
	wg.Wait()
	for v_i := 0; v_i < tempvideoi; v_i++ {
		writer.Write(images[v_i].Pix)
	}
	writer.Close()
	wg.Wait()
}

func writeDataToFrame(dat *os.File, i int, wg *sync.WaitGroup, images map[int]*image.RGBA, tempvideoi int) {
	data := make([]byte, bytesperframe*3)
	dat.ReadAt(data, int64(bytesperframe*3*i))
	image := encodeDataIntoPicture(data)
	mutexwrite.Lock()
	images[tempvideoi] = image
	mutexwrite.Unlock()
	wg.Done()
}

func main() {
	if len(os.Args) != 4 {
		printHelpScreen()
		return
	}
	if strings.ToLower(os.Args[1]) == "encode" {
		readFileWriteFrames(os.Args[2], os.Args[3])
		fmt.Println("finished")
		return
	}
	if strings.ToLower(os.Args[1]) == "decode" {
		readFrames(os.Args[3], os.Args[2])
		fmt.Println("finished")
		return
	}
	printHelpScreen()
}

func printHelpScreen() {
	fmt.Println("/goqrfile <encode/decode> <filename> <videofilename>")
}

func readFrames(videofilename string, outputfilename string) {
	resultfile, err := os.OpenFile(outputfilename, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println(err)
	}

	var bytetemp map[uint][]byte = make(map[uint][]byte)
	bytetempcounter := 0
	var wg sync.WaitGroup

	video, _ := vidio.NewVideo(videofilename)
	i := 0
	for video.Read() {
		wg.Add(1)
		t_image := image.NewRGBA(image.Rect(0, 0, video.Width(), video.Height()))
		copy(t_image.Pix, video.FrameBuffer())
		go getBytesFromFrame(t_image, &bytetemp, bytetempcounter, &wg)
		bytetempcounter++

		if i%1000 == 0 {
			fmt.Println(i)
			wg.Wait()
			for ii := 0; ii <= bytetempcounter; ii++ {
				resultfile.Write(bytetemp[uint(ii)])
			}

			bytetemp = make(map[uint][]byte)
			bytetempcounter = 0
		}
		i++
	}
	wg.Wait()
	for ii := 0; ii <= bytetempcounter; ii++ {
		resultfile.Write(bytetemp[uint(ii)])
	}
	resultfile.Close()
}

var mutex = &sync.Mutex{}

func getBytesFromFrame(src_image image.Image, bytetemp *map[uint][]byte, bytetempcounter int, wg *sync.WaitGroup) {
	var rgb_bytes map[int][]byte = make(map[int][]byte)
	rgb_bytes[0] = make([]byte, 0)
	rgb_bytes[1] = make([]byte, 0)
	rgb_bytes[2] = make([]byte, 0)

	n_row := 0
	n_col := 0

	for {
		var bits map[int]map[int]byte = make(map[int]map[int]byte)
		bits[0] = make(map[int]byte)
		bits[1] = make(map[int]byte)
		bits[2] = make(map[int]byte)
		for bit_n := 0; bit_n < 8; bit_n++ {
			if n_row > bitsPerRow-1 {
				n_row = 0
				n_col++
			}
			r, g, b, _ := src_image.At(n_row*bitPixelSize+(bitPixelSize/2), n_col*bitPixelSize+(bitPixelSize/2)).RGBA()
			if r > 0xffff/2 {
				bits[0][bit_n] = 1
			} else {
				bits[0][bit_n] = 0
			}
			if g > 0xffff/2 {
				bits[1][bit_n] = 1
			} else {
				bits[1][bit_n] = 0
			}
			if b > 0xffff/2 {
				bits[2][bit_n] = 1
			} else {
				bits[2][bit_n] = 0
			}
			n_row++
		}
		if n_col > bitsPerRow-1 {
			break
		}
		for i := 0; i < 3; i++ {
			t_byte := (bits[i][0] << 7) | (bits[i][1] << 6) | (bits[i][2] << 5) | (bits[i][3] << 4) | (bits[i][4] << 3) | (bits[i][5] << 2) | (bits[i][6] << 1) | bits[i][7]
			rgb_bytes[i] = append(rgb_bytes[i], t_byte)
		}

	}
	totalBytes := rgb_bytes[0]
	totalBytes = append(totalBytes, rgb_bytes[1]...)
	totalBytes = append(totalBytes, rgb_bytes[2]...)
	mutex.Lock()
	(*bytetemp)[uint(bytetempcounter)] = totalBytes
	mutex.Unlock()
	wg.Done()
}
