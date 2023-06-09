package main

import (
	"fmt"
	"image"
	"os"
	"strings"
	"time"

	"image/color"

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

func printInColor(color color.Color, character string) string {
	r, g, b, a := color.RGBA()
	//Scape to the rgba color in the terminal
	if a == 0 {
		//move the cursor the size of the character to the right

		return fmt.Sprintf("\033[%dC", len(character))
	}

	//set color of the character and background
	return fmt.Sprintf("\033[38;2;%d;%d;%dm\033[48;2;%d;%d;%dm%s\033[0m", r, g, b, r, g, b, character)
	//return fmt.Sprintf("\033[38;2;%d;%d;%dm%s\033[0m", r, g, b, character)
}

func printImage(img image.Image, fitX, fitY int) string {

	var buffer strings.Builder
	y := 0

	//save cursor position
	buffer.WriteString("\033[s")

	//print the image fitting the width and height
	for ; y < img.Bounds().Max.Y; y += fitY {
		for x := 0; x < img.Bounds().Max.X; x += fitX {
			buffer.WriteString(printInColor(img.At(x, y), "."))
		}
		buffer.WriteString("\n")
	}

	//restore cursor position
	buffer.WriteString("\033[u")

	//buffer.WriteString(fmt.Sprintf("\033[%dA", img.Bounds().Max.Y+1))

	//buffer.WriteString(fmt.Sprintf("\033[%dD", img.Bounds().Max.X+1))
	return buffer.String()

}

func generateGifCache(img *gif.GIF) []GIFCache {
	cache := make([]GIFCache, len(img.Image))
	scaleImage := img.Image[0]
	width := 150
	height := 200

	//keep the aspect ratio of the image
	if scaleImage.Bounds().Max.X > scaleImage.Bounds().Max.Y {
		height = int(float64(scaleImage.Bounds().Max.Y) / float64(scaleImage.Bounds().Max.X) * float64(width))
	} else {
		width = int(float64(scaleImage.Bounds().Max.X) / float64(scaleImage.Bounds().Max.Y) * float64(height))
	}

	addIterationX := scaleImage.Bounds().Max.X / width
	addIterationY := scaleImage.Bounds().Max.Y / height

	for i, frame := range img.Image {

		cache[i].images = printImage(frame, addIterationX, addIterationY)
		cache[i].delay = time.Duration(img.Delay[i]) * (time.Second / 100)
	}

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
