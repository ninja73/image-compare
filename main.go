package main

import (
	"fmt"
	"github.com/nfnt/resize"
	"image/jpeg"
	"bytes"
	"os"
	"log"
	"io/ioutil"
)

var BLOCK_SIZE uint = 20
func main() {
	fmt.Println("App runing...")
	img := NewImage("./files/miniony.jpg", "./files/smal-miniony.jpg")
	img.LoadImage()
}

func checkError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

type Image struct {
	fileName string
	outFileName string
}

func NewImage(fileName string, outFileName string) *Image {
	img := new(Image)
	img.fileName = fileName
	img.outFileName = outFileName
	return img
}

func(o *Image) LoadImage()  {
	fileLoad, fileError := ioutil.ReadFile(o.fileName)
	checkError(fileError)
	
	img, imgError := jpeg.Decode(bytes.NewReader(fileLoad))
	checkError(imgError)
	
	out, err := os.Create(o.outFileName)
	checkError(err)
	defer out.Close()
	
	small := resize.Resize(BLOCK_SIZE, BLOCK_SIZE, img, resize.Bilinear)
	bounds := small.Bounds()
	var t_data []int
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			r, g, b, _ := small.At(x, y).RGBA()
			t_data = append(t_data, (int(r) + int(g) + int(b)))
		}
	}
	fmt.Println(len(t_data))
	// write new image to file
	jpeg.Encode(out, small, nil)
}