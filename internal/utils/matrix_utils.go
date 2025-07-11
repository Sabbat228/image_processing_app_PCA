package utils

import (
	"errors"
	"image"
	"image/color"
	"math"

	_ "image/jpeg"

	"gonum.org/v1/gonum/mat"
)

var (
	ErrInvalidImageFormat = errors.New("invalid image format")
	ErrInvalidMethod      = errors.New("invalid processing method")
)

// ConvertImageToMatrix преобразует изображение в матрицу
func ConvertImageToMatrix(img image.Image) *mat.Dense {
	bounds := img.Bounds()
	width, height := bounds.Max.X-bounds.Min.X, bounds.Max.Y-bounds.Min.Y

	// Проверка размеров
	if width <= 0 || height <= 0 {
		return mat.NewDense(1, 1, []float64{0})
	}

	data := make([]float64, width*height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x+bounds.Min.X, y+bounds.Min.Y).RGBA()
			gray := 0.299*float64(r>>8) + 0.587*float64(g>>8) + 0.114*float64(b>>8)
			data[y*width+x] = gray
		}
	}

	return mat.NewDense(height, width, data)
}

// ConvertMatrixToImage преобразует матрицу в изображение
func ConvertMatrixToImage(m *mat.Dense) image.Image {
	rows, cols := m.Dims()
	img := image.NewGray(image.Rect(0, 0, cols, rows))

	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			var sum float64
			count := 0
			// Усреднение 3x3
			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= 1; dx++ {
					ny, nx := y+dy, x+dx
					if ny >= 0 && ny < rows && nx >= 0 && nx < cols {
						sum += m.At(ny, nx)
						count++
					}
				}
			}
			val := uint8(math.Round(sum / float64(count)))
			img.SetGray(x, y, color.Gray{Y: val})
		}
	}
	return img
}

// LoadImage загружает изображение из файла
