package main

import (
	"fmt"
	"image"
	"os"
	"strings"
	"time"

	"image/color"

	"golang.org/x/image/draw"

	"image/gif"
	_ "image/png"
)

// Images no included in this repository
const imageURL = "a.gif"

type GIFCache struct {
	images string
	delay  time.Duration
}

func main() {
	imagefile, err := os.Open(imageURL)
	if err != nil {
		panic(err)
	}

	gifImage, err := gif.DecodeAll(imagefile)

	if err != nil {
		panic(err)
	}

	cache := generateGifCache(gifImage)
	for {
		printGif(cache)
	}
}

func printInColor(color color.Color) string {
	r, g, b, a := color.RGBA()
	//Scape to the rgba color in the terminal
	if a == 0 {
		//move the cursor 1 position to the right
		return "\033[1C"
	}

	//same color for the character and the background
	if r == g && g == b {
		//move the cursor 1 position to the right
		return "\033[1C"
	}

	return fmt.Sprintf("\033[38;2;%d;%d;%dm\033[48;2;%d;%d;%dm█\033[0m", r, g, b, r, g, b)
}

func printImage(img image.Image) string {
	var buffer strings.Builder
	y := 0
	countNewLine := 1

	//width := img.Bounds().Max.X / 2
	//height := img.Bounds().Max.Y / 2

	width := 100
	height := 100

	img = resizeImage(img, width, height)
	for ; y < img.Bounds().Max.Y; y += 1 {
		for x := 0; x < img.Bounds().Max.X; x += 1 {
			buffer.WriteString(printInColor(img.At(x, y)))
			buffer.WriteString(printInColor(img.At(x, y)))
			buffer.WriteString(printInColor(img.At(x, y)))
		}

		buffer.WriteString("\n")
		countNewLine++
	}

	buffer.WriteString(fmt.Sprintf("\033[%dA", countNewLine))

	buffer.WriteString(fmt.Sprintf("\033[%dD", img.Bounds().Max.X))
	return buffer.String()

}

func resizeImage(src image.Image, width, height int) image.Image {
	img := image.NewRGBA(
		image.Rect(0, 0, width, height),
	)
	draw.NearestNeighbor.Scale(img, img.Bounds(), src, src.Bounds(), draw.Over, nil)
	return img
}

func generateGifCache(img *gif.GIF) []GIFCache {
	cache := make([]GIFCache, len(img.Image))
	for i, frame := range img.Image {

		cache[i].images = printImage(frame)
		cache[i].delay = time.Duration(img.Delay[i]) * (time.Second / 100)
	}

	return cache
}

func printGif(cache []GIFCache) {
	for _, frame := range cache {
		fmt.Print(frame.images)
		time.Sleep(frame.delay)
	}
}
