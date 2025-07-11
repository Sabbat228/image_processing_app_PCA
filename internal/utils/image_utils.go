package utils

import (
	"image"
	"os"
	"path/filepath"

	_ "image/jpeg"
	"image/png"
)

// LoadImage загружает изображение из файла
func LoadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// SaveImage сохраняет изображение в файл (формат PNG)
func SaveImage(path string, img image.Image) error {
	// Создаем директорию если ее нет
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, img)
}
