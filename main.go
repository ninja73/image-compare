package main

import (
	"fmt"
	"github.com/nfnt/resize"
	"image/jpeg"
	"bytes"
	"log"
	"io/ioutil"
	"strings"
	"os"
	"math"
	"strconv"
)

const (
	BLOCK_SIZE uint = 20
	THRESHOLD int = 60
)

func main() {
	fmt.Println("App runing...")
	imageList := NewImageList("./files/")
	
	for _, img := range imageList.images {
		result := []string{}
		for _, x := range imageList.images {
			distance := img.Compare(x)
			if distance > 220 {
				result = append(result, x.filename + ": " + strconv.Itoa(distance))
			}
		}
		fmt.Println(img.filename+ ":" + strings.Join(result, ""))
	}
	
}

type ResultImg struct {
	Image string
	Distances string
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
			o.t_data = append(o.t_data, int(uint8(r) + uint8(g) + uint8(b)))
		}
	}
	
	return o
}

func(o *Image) Compare(other Image) int {
	var result int = 0
	data := difference(o.t_data, other.t_data)
	for _, x := range data {
		if(x > THRESHOLD) {
			result += x
		}
	}
	return result
}

func difference(data1 []int, data2 []int) []int {
	var diff []int
	if len(data1) == len(data2) {
		for i := 0; i < len(data1); i++ {
			t := math.Abs(float64(data1[i] - data2[i]))
			
			diff = append(diff, int(t))
		}
	}
	return diff
}

/** Подгружаем файлы из директории
Читаем файлы формата *.jpg
 */
type ImageList struct {
	dirname string
	images [] Image
}

func NewImageList(pathname string) *ImageList {
	imageList := new(ImageList)
	imageList.dirname = pathname
	imageList.Load()
	return imageList
}

func(o *ImageList) Load () *ImageList  {
	files, _ := ioutil.ReadDir(o.dirname)
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".jpg") {
			fullPath := o.dirname + string(os.PathSeparator) + file.Name()
			o.images = append(o.images, *NewImage(fullPath).Load())
		}
	}
	
	return o
}
