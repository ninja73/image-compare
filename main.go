package main

import (
	"fmt"
	"github.com/nfnt/resize"
	"image/jpeg"
	"bytes"
	"log"
	"io/ioutil"
	"strings"
	//"os"
)

var BLOCK_SIZE uint = 20
func main() {
	fmt.Println("App runing...")
	NewImageList("./files/")
}

func checkError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
/**
Получаем информацию о файле
 */
type Image struct {
	filename string
	t_data [] int
}

func NewImage(fileName string) *Image {
	img := new(Image)
	img.filename = fileName
	return img
}

func(o *Image) Load() *Image {
	fileLoad, fileError := ioutil.ReadFile(o.filename)
	checkError(fileError)
	
	img, imgError := jpeg.Decode(bytes.NewReader(fileLoad))
	checkError(imgError)
	
	small := resize.Resize(BLOCK_SIZE, BLOCK_SIZE, img, resize.Bilinear)
	bounds := small.Bounds()
	
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			r, g, b, _ := small.At(x, y).RGBA()
			o.t_data = append(o.t_data, (int(r) + int(g) + int(b)))
		}
	}
	return o
}

/** Подгружаем файлы из директории
Читаем файлы формата *.jpg
 */
type ImageList struct {
	dirname string
	images [] Image
}

func NewImageList(pathname string) ImageList {
	imageList := new(ImageList)
	imageList.dirname = pathname
	imageList.Load()
	return imageList
}

func(o *ImageList) Load () *ImageList  {
	files, _ := ioutil.ReadDir(o.dirname)
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".jpg") {
			o.images = append(o.images, NewImage(file.Name()).Load())
		}
	}
	fmt.Println(o.images)
	return o
}
