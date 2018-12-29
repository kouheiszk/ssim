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

func convertToGrayscale(img image.Image) image.Image {
	imageBounds := img.Bounds()
	bounds := img.Bounds()
	dest := image.NewGray(imageBounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			originalColor := img.At(x, y)
			grayColor := color.GrayModel.Convert(originalColor)
			dest.Set(x, y, grayColor)
		}
	}

	return dest
}

func pixelValue(c color.Color) float64 {
	r, _, _, _ := c.RGBA()
	return float64(r >> 8)
}

func isSameDimension(img1, img2 image.Image) bool {
	bounds1 := img1.Bounds()
	bounds2 := img2.Bounds()
	return (bounds1.Dx() == bounds2.Dx()) && (bounds1.Dy() == bounds2.Dy())
}

func average(img image.Image, bounds image.Rectangle) float64 {
	sum := 0.0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			sum += pixelValue(img.At(x, y))
		}
	}

	pixels := float64(bounds.Dx() * bounds.Dy())
	return sum / pixels
}

func variance(img image.Image, bounds image.Rectangle, avg float64) float64 {
	sum := 0.0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			value := pixelValue(img.At(x, y))
			sum += math.Pow(value-avg, 2.0)
		}
	}

	pixels := float64(bounds.Dx() * bounds.Dy())
	return sum / pixels
}

func covariance(img1, img2 image.Image, bounds image.Rectangle, avg1, avg2 float64) float64 {
	sum := 0.0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			value1 := pixelValue(img1.At(x, y))
			value2 := pixelValue(img2.At(x, y))
			sum += (value1 - avg1) * (value2 - avg2)
		}
	}

	pixels := float64(bounds.Dx() * bounds.Dy())
	return sum / pixels
}

func ssim(img1, img2 image.Image, bounds image.Rectangle) float64 {
	avg1 := average(img1, bounds)
	avg2 := average(img2, bounds)

	var1 := variance(img1, bounds, avg1)
	var2 := variance(img2, bounds, avg2)

	covar := covariance(img1, img2, bounds, avg1, avg2)

	num := ((2.0 * avg1 * avg2) + C1) * ((2.0 * covar) + C2)
	denominator := (math.Pow(avg1, 2.0) + math.Pow(avg2, 2.0) + C1) * (var1 + var2 + C2)
	return num / denominator
}

func SSIM(img1, img2 image.Image) (float64, error) {
	if !isSameDimension(img1, img2) {
		return 0, fmt.Errorf("images must have the same dimension")
	}

	img1 = convertToGrayscale(img1)
	img2 = convertToGrayscale(img2)

	bounds := img1.Bounds()
	windowSizeX := bounds.Dx() / 4
	windowSizeY := bounds.Dy() / 4

	sum := 0.0
	windows := 0
	for y := 0; y <= bounds.Dy()-windowSizeY; y += windowSizeY / 2 {
		for x := 0; x <= bounds.Dx()-windowSizeX; x += windowSizeX / 2 {
			window := image.Rect(x, y, x+windowSizeX, y+windowSizeY)
			sum += ssim(img1, img2, window)
			windows++
		}
	}

	return sum / float64(windows), nil
}
