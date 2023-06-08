package main

import (
	"fmt"
	"image"
	"os"
	"time"

	"image/color"
	"image/gif"
	_ "image/png"
)

const imageURL = "a.gif"

func main() {
	imagefile, err := os.Open(imageURL)
	if err != nil {
		panic(err)
	}

	gifImage, err := gif.DecodeAll(imagefile)

	if err != nil {
		panic(err)
	}

	printGiftImage(gifImage)
}

func printInColor(color color.Color) string {
	r, g, b, a := color.RGBA()
	//Scape to the rgba color in the terminal
	if a == 0 {
		//move the cursor 1 position to the right
		return "\033[1C"
	}

	return fmt.Sprintf("\033[38;2;%d;%d;%dm█\033[0m", r, g, b)
}

func printImage(img image.Image) {
	y := 0
	countNewLine := 1
	for ; y < img.Bounds().Max.Y; y += 3 {
		for x := 0; x < img.Bounds().Max.X; x += 3 {
			fmt.Print(printInColor(img.At(x, y)))
			fmt.Print(".")
		}
		fmt.Println()
		countNewLine++
	}

	//move the cursor up
	fmt.Printf("\033[%dA", countNewLine)

	//move the cursor to the left beginning of the line
	fmt.Printf("\033[%dD", img.Bounds().Max.X/20)
}

func printGiftImage(img *gif.GIF) {
	for i, frame := range img.Image {
		printImage(frame)
		//delay
		fmt.Println(img.Delay[i])

		time.Sleep(time.Duration(img.Delay[i]) * time.Microsecond)
	}
	printGiftImage(img)
}
