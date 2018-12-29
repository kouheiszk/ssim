package main

import (
	"fmt"
	"github.com/kouheiszk/ssim"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

func main() {
	path1 := "../../test/pic.jpg"
	path2 := "../../test/pic2.jpg"

	img1, err := loadImage(path1)
	if err != nil {
		panic(err)
	}

	img2, err := loadImage(path2)
	if err != nil {
		panic(err)
	}

	ssim, err := ssim.SSIM(img1, img2)
	if err != nil {
		panic(err)
	}

	fmt.Printf("SSIM = %f", ssim)
}

func loadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	return img, err
}
