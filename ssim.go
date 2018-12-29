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

const (
	R = 0
	G = 1
	B = 2
)

func pixelValue(c color.Color, component int) float64 {
	r, g, b, _ := c.RGBA()
	return float64([]uint32{r, g, b}[component] >> 8)
}

func isSameDimension(img1, img2 image.Image) bool {
	bounds1 := img1.Bounds()
	bounds2 := img2.Bounds()
	return (bounds1.Dx() == bounds2.Dx()) && (bounds1.Dy() == bounds2.Dy())
}

func average(img image.Image, component int) float64 {
	bounds := img.Bounds()
	sum := 0.0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			sum += pixelValue(img.At(x, y), component)
		}
	}

	pixels := float64(bounds.Dx() * bounds.Dy())
	return sum / pixels
}

func variance(img image.Image, avg float64, component int) float64 {
	bounds := img.Bounds()
	sum := 0.0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			value := pixelValue(img.At(x, y), component)
			sum += math.Pow(value-avg, 2.0)
		}
	}

	pixels := float64(bounds.Dx() * bounds.Dy())
	return sum / pixels
}

func covariance(img1, img2 image.Image, avg1, avg2 float64, component int) float64 {
	bounds := img1.Bounds()
	sum := 0.0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			value1 := pixelValue(img1.At(x, y), component)
			value2 := pixelValue(img2.At(x, y), component)
			sum += (value1 - avg1) * (value2 - avg2)
		}
	}

	pixels := float64(bounds.Dx() * bounds.Dy())
	return sum / pixels
}

func ssim(img1, img2 image.Image, component int) float64 {
	avg1 := average(img1, component)
	avg2 := average(img2, component)

	var1 := variance(img1, avg1, component)
	var2 := variance(img2, avg2, component)

	covar := covariance(img1, img2, avg1, avg2, component)

	num := ((2.0 * avg1 * avg2) + C1) * ((2.0 * covar) + C2)
	denominator := (math.Pow(avg1, 2.0) + math.Pow(avg2, 2.0) + C1) * (var1 + var2 + C2)
	return num / denominator
}

func SSIM(img1, img2 image.Image) (float64, error) {
	if !isSameDimension(img1, img2) {
		return 0, fmt.Errorf("images must have the same dimension")
	}

	r := ssim(img1, img2, R)
	g := ssim(img1, img2, G)
	b := ssim(img1, img2, B)

	return (r + g + b) / 3.0, nil
}
