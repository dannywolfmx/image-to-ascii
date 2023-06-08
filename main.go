package main

import (
	"fmt"
	"image"
	"os"

	"image/color"
	_ "image/png"
)

const imageURL = "pato.png"

func main() {
	imagefile, err := os.Open(imageURL)
	if err != nil {
		panic(err)
	}

	img, format, err := image.Decode(imagefile)

	fmt.Println(format)

	if err != nil {
		panic(err)
	}

	buffer := ""
	for y := 0; y < img.Bounds().Max.Y; y += 2 {
		for x := 0; x < img.Bounds().Max.X; x += 2 {

			buffer += printInColor(img.At(x, y))
		}
		buffer += "\n"
	}

	fmt.Println(buffer)
}

func printInColor(color color.Color) string {
	r, g, b, _ := color.RGBA()
	//Scape to the rgba color in the terminal
	return fmt.Sprintf("\033[38;2;%d;%d;%dm█\033[0m", r, g, b)
}
