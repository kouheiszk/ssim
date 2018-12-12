package main

import (
	"fmt"
	"github.com/kouheiszk/ssim"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"reflect"
)

func main() {
	file, err := os.Open("../../test/pic.jpg")
	if err != nil {
		log.Fatalf("Could Not Open Pic -> %v", err)
	}

	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatalf("Could Not Decode Pic -> %v", err)
	}

	file, err = os.Open("../../test/pic2.jpg")
	if err != nil {
		log.Fatalf("Could Not Open Pic2 -> %v", err)
	}

	img2, _, err := image.Decode(file)
	if err != nil {
		log.Fatalf("Could Not Decode Pic2 -> %v", err)
	}
	index, err := ssim.SSIM(img, img2)
	if err != nil {
		log.Fatalf("Could Not Calculate SSIM -> %v", err)
	}

	fmt.Printf("Index is type of = %v\n", reflect.TypeOf(index))
	fmt.Printf("Index Value = %f", index)
}
