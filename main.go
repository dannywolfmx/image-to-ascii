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
const imageURL = "cow.gif"

func main() {
	imagefile, err := os.Open(imageURL)
	if err != nil {
		panic(err)
	}

	gifImage, err := gif.DecodeAll(imagefile)

	if err != nil {
		panic(err)
	}

	for {
		printGiftImage(gifImage)
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

func printImage(img image.Image) {
	var buffer strings.Builder
	y := 0
	countNewLine := 1
	for ; y < img.Bounds().Max.Y; y += 3 {
		for x := 0; x < img.Bounds().Max.X; x += 1 {
			buffer.WriteString(printInColor(img.At(x, y)))
			//square character
		}

		buffer.WriteString("\n")
		countNewLine++
	}
	fmt.Println(buffer.String())

	//move the cursor up
	fmt.Printf("\033[%dA", countNewLine+1)

	//move the cursor to the left beginning of the line
	fmt.Printf("\033[%dD", img.Bounds().Max.X/3)

}

func printGiftImage(img *gif.GIF) {
	for i, frame := range img.Image {
		printImage(frame)
		//delay
		fmt.Println(img.Delay[i])

		time.Sleep(time.Duration(img.Delay[i]) * time.Microsecond)
	}
}
