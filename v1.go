package main

import (
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"math"
	"os"
	"strconv"
	"sync"

	vidio "github.com/AlexEidt/Vidio"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/datamatrix"
)

var readamounts int = 512
var imagesize int = 1024

func main() {
	//createFramesFromFile()
	readFiles()

}

func readFiles() {
	var totalouputbytes map[int][]byte = make(map[int][]byte)
	var wg sync.WaitGroup
	for i := 0; true; i++ {
		number := strconv.Itoa(i + 1)
		paddinglength := 5 - len(number)
		for iiii := 0; iiii < paddinglength; iiii++ {
			number = "0" + number
		}
		//fmt.Println(number)
		f, err := os.Open("images/image" + number + ".png")
		if err != nil {
			break
		}
		wg.Add(1)
		go readIndividualFrame(i*3, &totalouputbytes, &wg, f)
		if i%50 == 0 {
			fmt.Println("i:", i, " waiting")
			wg.Wait()
		}
	}
	wg.Wait()
	resultfile, err := os.OpenFile("result.gif", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println(err)
	}
	var resulttotalbytesoutput []byte
	for i := 0; i < len(totalouputbytes); i++ {
		resulttotalbytesoutput = append(resulttotalbytesoutput, totalouputbytes[i]...)
	}
	resultfile.Write(resulttotalbytesoutput)
	resultfile.Close()
}

var mutex = &sync.Mutex{}

func readIndividualFrame(i int, outputarray *map[int][]byte, wg *sync.WaitGroup, f *os.File) {
	defer wg.Done()
	defer f.Close()

	var decodedQrMap map[int]map[int]map[int]uint32 = make(map[int]map[int]map[int]uint32)
	decodedQrMap[0] = make(map[int]map[int]uint32)
	decodedQrMap[1] = make(map[int]map[int]uint32)
	decodedQrMap[2] = make(map[int]map[int]uint32)

	src_image, _, err := image.Decode(f)
	if err != nil {
		fmt.Println(err)
	}
	for x := 0; x < imagesize; x++ {
		decodedQrMap[0][x] = make(map[int]uint32)
		decodedQrMap[1][x] = make(map[int]uint32)
		decodedQrMap[2][x] = make(map[int]uint32)
		for y := 0; y < imagesize; y++ {
			r, g, b, _ := src_image.At(x, y).RGBA()
			decodedQrMap[0][x][y] = r
			decodedQrMap[1][x][y] = g
			decodedQrMap[2][x][y] = b
		}
	}

	for ii := 0; ii < 3; ii++ {
		resultimg := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{imagesize, imagesize}})
		for x := 0; x < imagesize; x++ {
			for y := 0; y < imagesize; y++ {
				resultimg.Set(x, y, color.RGBA{uint8(decodedQrMap[ii][x][y]), uint8(decodedQrMap[ii][x][y]), uint8(decodedQrMap[ii][x][y]), 0xff})
			}
		}
		reader := datamatrix.NewDataMatrixReader()
		bmp, _ := gozxing.NewBinaryBitmapFromImage(resultimg)
		result, err := reader.Decode(bmp, nil)
		if err == nil {
			qrqrtext := result.GetText()
			qrqrbytes, _ := base64.RawStdEncoding.DecodeString(qrqrtext)
			//fmt.Println("seting:", i)
			mutex.Lock()
			(*outputarray)[i+ii] = qrqrbytes
			mutex.Unlock()
		}
	}
}

func createFramesFromFile() {
	dat, _ := os.Open("source.gif")
	defer dat.Close()
	fileinfo, _ := dat.Stat()
	readcycles := math.Ceil((float64(fileinfo.Size()) / float64(readamounts)) / 3)
	options := vidio.Options{
		FPS: 25,
	}
	writer, _ := vidio.NewVideoWriter("output.mp4", imagesize, imagesize, &options)
	for x := 0; x < int(math.Ceil(readcycles/10)); x++ {
		var wg sync.WaitGroup
		var images map[int]*image.RGBA = make(map[int]*image.RGBA)
		tempi := 0
		for i := x * 20; i < (x*20)+20 && i < int(readcycles); i++ {
			wg.Add(1)
			go createQrAndWriteFrame(i, dat, readcycles, &wg, images, tempi)
			tempi++
		}
		wg.Wait()

		for i := 0; i < len(images); i++ {
			writer.Write(images[i].Pix)
		}
	}
	writer.Close()
}

var mutex2 = &sync.Mutex{}

func createQrAndWriteFrame(i int, dat *os.File, readcycles float64, wg *sync.WaitGroup, imagemap map[int]*image.RGBA, tempi int) {
	defer wg.Done()
	if i%100 == 0 {
		fmt.Println(math.Ceil((float64(i)/readcycles)*100), "%")
	}
	var images []image.Image = make([]image.Image, 3)
	for ii := 0; ii < 3; ii++ {
		red_channel := make([]byte, readamounts)
		dat.ReadAt(red_channel, int64(readamounts*i*3+(readamounts*ii)))

		writer := datamatrix.NewDataMatrixWriter()
		png, err := writer.EncodeWithoutHint(base64.RawStdEncoding.EncodeToString(red_channel), gozxing.BarcodeFormat_DATA_MATRIX, imagesize, imagesize)
		if err != nil {
			panic(err)
		}
		images[ii] = png
	}
	var rgbvalues map[int]map[int]map[int]uint8 = make(map[int]map[int]map[int]uint8)
	for ii := 0; ii < 3; ii++ {
		if images[ii] == nil {
			fmt.Println("break")
			break
		}
		rgbvalues[ii] = make(map[int]map[int]uint8)
		img := images[ii]
		for x := 0; x < imagesize; x++ {
			rgbvalues[ii][x] = make(map[int]uint8)
			for y := 0; y < imagesize; y++ {
				r, _, _, _ := img.At(x, y).RGBA()
				if r > 0 {
					rgbvalues[ii][x][y] = 0xff
				} else {
					rgbvalues[ii][x][y] = 0
				}
			}
		}
	}
	resultimg := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{imagesize, imagesize}})
	for x := 0; x < imagesize; x++ {
		for y := 0; y < imagesize; y++ {
			resultimg.Set(x, y, color.RGBA{uint8(rgbvalues[0][x][y]), uint8(rgbvalues[1][x][y]), uint8(rgbvalues[2][x][y]), 0xff})
		}
	}
	mutex2.Lock()
	imagemap[tempi] = resultimg
	mutex2.Unlock()
}
