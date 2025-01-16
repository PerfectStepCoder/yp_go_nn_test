package httpp

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"mime/multipart"
)

// ImageToFloat32Matrix - функция для нормализации пикселей в диапазоне [-1, 1]
func ImageToFloat32Matrix(file multipart.File) ([][]float32, error) {
	// Декодируем изображение из файла
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("ошибка декодирования изображения: %w", err)
	}

	// Преобразуем изображение в grayscale (или используем цветное, если нужно)
	bounds := img.Bounds()
	grayImg := image.NewGray(bounds) // Создаём новое grayscale изображение

	width, height := bounds.Dx(), bounds.Dy()
	matrix := make([][]float32, height)

	for y := 0; y < height; y++ {
		matrix[y] = make([]float32, width)
		for x := 0; x < width; x++ {
			// Получаем цвет пикселя
			pixel := grayImg.At(x, y)

			// Преобразуем в яркость (grayscale)
			gray := color.GrayModel.Convert(pixel).(color.Gray)

			// Нормализуем в диапазон [-1, 1]
			normalized := float32(gray.Y)/127.5 - 1

			// Заполняем матрицу
			matrix[y][x] = normalized
		}
	}

	return matrix, nil
}

func getSize(input [][]float32) (height int, width int) {
	// Вычисляем высоту и ширину
	height = len(input) // Количество строк
	width = 0
	if height > 0 {
		width = len(input[0]) // Количество столбцов в первой строке
	}
	return height, width
}

// Функция для преобразования [A][B]float32 в [1][AxB]float32
func ReshapeTo1xN(input [][]float32) [][]float32 {

	height, width := getSize(input)

	output := make([][]float32, 1)
	output[0] = make([]float32, height*width)

	index := 0

	// Копируем элементы построчно
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			output[0][index] = input[i][j]
			index++
		}
	}

	return output
}
