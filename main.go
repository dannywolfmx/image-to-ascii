package main

import (
	"flag"
	"fmt"
	"image"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"image/color"

	"image/gif"
	_ "image/png"
)

type GIFCache struct {
	images string
	delay  time.Duration
}

//go run main.go -h=75 -w=100 -i=https://media.giphy.com/media/3o7aDcz2YkI2XK8MlO/giphy.gif

func main() {
	height := flag.Int("h", 25, "height of the image")
	width := flag.Int("w", 25, "width of the image")
	imageURL := flag.String("i", "", "image url")

	flag.Parse()

	if *imageURL == "" {
		fmt.Println("===== You need to specify an image url =====")
		return
	}

	var imageFile io.ReadCloser

	//check if the image is a local file or a url
	if strings.HasPrefix(*imageURL, "http") {
		imageFile = getImageFromURL(*imageURL)
	} else {
		_, err := os.Stat(*imageURL)
		if err != nil {
			fmt.Println("===== The image file does not exist =====")
			return
		}
		imageFile, err = os.Open(*imageURL)
		if err != nil {
			panic(err)
		}
	}

	gifImage, err := gif.DecodeAll(imageFile)
	imageFile.Close()

	if err != nil {
		panic(err)
	}

	cache := generateGifCache(gifImage, *height, *width)
	for {
		printGif(cache)
	}
}

func printInColor(w io.Writer, color color.Color, character string) {
	r, g, b, a := color.RGBA()
	//Scape to the rgba color in the terminal
	if a == 0 {
		//move the cursor the size of the character to the right
		fmt.Fprintf(w, "\033[%dC", len(character))
		return
	}

	//set color of the character and background
	//return fmt.Sprintf("\033[38;2;%d;%d;%dm\033[48;2;%d;%d;%dm", r, g, b, r, g, b)

	fmt.Fprintf(w, "\033[38;2;%d;%d;%dm%s\033[0m", r, g, b, character)
}

func getImageFromURL(url string) io.ReadCloser {
	response, err := http.Get(url)

	if err != nil {
		panic(err)
	}

	return response.Body
}

func printImage(img image.Image, fitX, fitY int) string {

	var buffer strings.Builder
	y := 0

	//save cursor position
	fmt.Fprint(&buffer, "\033[s")

	//transform to grayscale

	grayImage := image.NewGray16(img.Bounds())

	for y := 0; y < img.Bounds().Max.Y; y++ {
		for x := 0; x < img.Bounds().Max.X; x++ {
			grayImage.Set(x, y, img.At(x, y))
		}
	}

	//print the image fitting the width and height
	for ; y < img.Bounds().Max.Y; y += fitY {
		for x := 0; x < img.Bounds().Max.X; x += fitX {
			//grayScale
			grayImage.Set(x, y, img.At(x, y))
			r, g, b, _ := grayImage.At(x, y).RGBA()
			//Note the gray convertion will delete the alpha channel,so we need to
			//get the alpha channel from the original image
			_, _, _, a := img.At(x, y).RGBA()
			if a == 0 {
				//move the cursor the size of the character to the right
				fmt.Fprintf(&buffer, "\033[1C")
				continue
			}

			//printInColor(&buffer, grayImage.At(x, y), "#")

			//	//set color as ascii
			fmt.Fprint(&buffer, string(gray16ToAnsi(r, g, b)))

		}
		fmt.Fprintln(&buffer, "")
	}

	//restore cursor position
	fmt.Fprint(&buffer, "\033[u")

	//buffer.WriteString(fmt.Sprintf("\033[%dA", img.Bounds().Max.Y+1))

	//buffer.WriteString(fmt.Sprintf("\033[%dD", img.Bounds().Max.X+1))
	return buffer.String()

}

func generateGifCache(img *gif.GIF, height, width int) []GIFCache {
	cache := make([]GIFCache, len(img.Image))
	scaleImage := img.Image[0]

	addIterationX := scaleImage.Bounds().Max.X / width
	addIterationY := scaleImage.Bounds().Max.Y / height

	var wg sync.WaitGroup

	for i, frame := range img.Image {
		wg.Add(1)
		go (func(index int, paletted image.Image) {
			defer wg.Done()

			cache[index].images = printImage(paletted, addIterationX, addIterationY)
			cache[index].delay = time.Duration(img.Delay[index]) * (time.Second / 100)
		})(i, frame)
	}

	wg.Wait()

	return cache
}

func printGif(cache []GIFCache) {
	//clear screen
	fmt.Print("\033[H\033[2J")
	for _, frame := range cache {
		fmt.Print(frame.images)
		time.Sleep(frame.delay)
	}
}

// Transform RGB to ANSII
func gray16ToAnsi(r, g, b uint32) rune {
	r8, _, _ := rgb16ToRgb8(r, g, b)

	if r8 >= 192 {
		return '█'
	}

	if r8 >= 128 {
		return '▓'
	}

	if r8 >= 64 {
		return '▒'
	}
	return ' '

}

// Note the gray color is the same for all the colors example
// if r is 255 the g and b will be the same in all the cases
func rgb16ToRgb8(r, g, b uint32) (uint8, uint8, uint8) {
	gray8Value := uint8(r >> 8)
	return gray8Value, gray8Value, gray8Value
}
