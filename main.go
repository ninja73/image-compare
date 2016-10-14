package main

import (
	"fmt"
	"image/jpeg"
	"bytes"
	"log"
	"io/ioutil"
	"strings"
	"os"
	"math"
	"strconv"
	"github.com/disintegration/gift"
	"image"
)

const (
	BLOCK_SIZE int = 20
	THRESHOLD int = 60
)

func main() {
	fmt.Println("App runing...")
	imageList := NewImageList("./files/")
	
	for _, img := range imageList.images {
		result := []string{}
		for _, x := range imageList.images {
			distance := img.Compare(x)
			if distance < 220 {
				result = append(result, "\r\n"+ "path: " + x.path +
								", filename: " + x.filename + ": " + strconv.Itoa(distance))
			}
		}
		fmt.Println("--------" + img.filename + "--------" + strings.Join(result, ""))
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
	path, filename string
	t_data [] int
}

func NewImage(path string, fileName string) *Image {
	img := new(Image)
	img.filename = fileName
	img.path = path
	return img
}

func(o *Image) Load() *Image {
	fileLoad, fileError := ioutil.ReadFile(o.path + string(os.PathSeparator) + o.filename)
	checkError(fileError)
	
	img, imgError := jpeg.Decode(bytes.NewReader(fileLoad))
	checkError(imgError)
	
	dstImage := image.NewRGBA(image.NewRGBA(image.Rect(0, 0, BLOCK_SIZE, BLOCK_SIZE)).Bounds())
	
	small := gift.New(
		gift.ResizeToFill(BLOCK_SIZE, BLOCK_SIZE, gift.LanczosResampling, gift.CenterAnchor),
		gift.Grayscale(),
		gift.Sobel())
	
	small.Draw(dstImage, img)
	
	out, err := os.Create("tmp-"+o.filename)
	checkError(err)
	defer out.Close()
	jpeg.Encode(out, dstImage, nil)
	
	bounds := dstImage.Bounds()
	
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			r, g, b, _ := dstImage.At(x, y).RGBA()
			o.t_data = append(o.t_data, int((uint8(r) + uint8(g) + uint8(b))))
		}
	}
	
	return o
}

/** Находим растояние между картинками
 */
func(o *Image) Compare(other Image) int {
	return difference(o.t_data, other.t_data)
}

func difference(data1 []int, data2 []int) int {
	result := 0
	
	if len(data1) == len(data2) {
		for i := 0; i < len(data1); i++ {
			t := int(math.Abs(float64(data1[i] - data2[i])))
			if(t > THRESHOLD) {
				result++
			}
		}
	}
	return result
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
			o.images = append(o.images, *NewImage(o.dirname, file.Name()).Load())
		}
	}
	
	return o
}
/*out, err := os.Create(o.outFileName)
	checkError(err)
	defer out.Close()*/

// write new image to file
//jpeg.Encode(out, small, nil)