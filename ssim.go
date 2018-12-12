package ssim

import (
	"fmt"
	"image"
	"image/color"
	"math"
)

// Default SSIM constants
var (
	L  = 255.0
	K1 = 0.01
	K2 = 0.03
	C1 = math.Pow(K1*L, 2.0)
	C2 = math.Pow(K2*L, 2.0)
)

func dimensions(img image.Image) (width, height int) {
	bounds := img.Bounds()
	return bounds.Dx(), bounds.Dy()
}

func getPixelValue(c color.Color) float64 {
	r, _, _, _ := c.RGBA()
	return float64(r >> 8)
}

func equalDim(img1, img2 image.Image) bool {
	w1, h1 := dimensions(img1)
	w2, h2 := dimensions(img2)
	return (w1 == w2) && (h1 == h2)
}

func convertToGrayscale(img image.Image) image.Image {
	imageBounds := img.Bounds()
	width, height := dimensions(img)

	grayImage := image.NewGray(imageBounds)

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			originalColor := img.At(x, y)
			grayColor := color.GrayModel.Convert(originalColor)
			grayImage.Set(x, y, grayColor)
		}
	}

	return grayImage
}

func mean(img image.Image) float64 {
	width, height := dimensions(img)
	n := float64((width * height) - 1)
	sum := 0.0
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			sum += getPixelValue(img.At(x, y))
		}
	}

	return sum / n
}

func stdDev(img image.Image) float64 {
	width, height := dimensions(img)

	n := float64((width * height) - 1)
	sum := 0.0
	avg := mean(img)

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			pixel := getPixelValue(img.At(x, y))
			sum += math.Pow(pixel-avg, 2.0)
		}
	}

	return math.Sqrt(sum / n)
}

func covar(img1, img2 image.Image) (float64, error) {
	if !equalDim(img1, img2) {
		return 0, fmt.Errorf("images must have the same dimension")
	}

	avg1 := mean(img1)
	avg2 := mean(img2)

	width, height := dimensions(img1)
	sum := 0.0
	n := float64((width * height) - 1)

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			pix1 := getPixelValue(img1.At(x, y))
			pix2 := getPixelValue(img2.At(x, y))
			sum += (pix1 - avg1) * (pix2 - avg2)
		}
	}

	return sum / n, nil
}

func SSIM(x, y image.Image) (float64, error) {
	x = convertToGrayscale(x)
	y = convertToGrayscale(y)

	avgX := mean(x)
	avgY := mean(y)

	stdDevX := stdDev(x)
	stdDevY := stdDev(y)

	covar, err := covar(x, y)
	if err != nil {
		return 0, err
	}

	num := ((2.0 * avgX * avgY) + C1) * ((2.0 * covar) + C2)
	denominator := (math.Pow(avgX, 2.0) + math.Pow(avgY, 2.0) + C1) * (math.Pow(stdDevX, 2.0) + math.Pow(stdDevY, 2.0) + C2)
	return num / denominator, nil
}
