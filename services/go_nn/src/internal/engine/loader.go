package engine

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

// LoadDataset - загрузка датасета MNIST (размер изображений 28x28)
func LoadDataset(filename string, batchSize int) ([][][]float32, [][]int, error) {  // old ([][784]float32, []int, error) 
	// Открываем файл
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	// Читаем содержимое CSV
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, nil, err
	}

	// Преобразуем данные
	var images [][784]float32
	var labels []int
	for _, record := range records[1:] { // Пропускаем заголовок
		// Первая колонка — метка
		label, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, nil, err
		}
		labels = append(labels, label)

		// Остальные колонки — пиксели
		var image [784]float32
		for i := 0; i < 784; i++ {
			pixel, err := strconv.ParseFloat(record[i+1], 32)
			if err != nil {
				return nil, nil, err
			}
			image[i] = float32(pixel)
		}
		images = append(images, image)
	}
	//fmt.Println("All images:", len(images))
	
	return batchArrayImages(images, batchSize), batchArrayLabels(labels, batchSize), nil
}


// batchArrayImages - функция для нарезки на батчи
func batchArrayImages(data [][784]float32, batchSize int) [][][]float32 {
	if batchSize <= 0 {
		panic("batchSize должен быть больше 0")
	}

	var batches [][][]float32
	for i := 0; i < len(data); i += batchSize {
		// Определяем конец текущего батча
		end := i + batchSize
		if end > len(data) {
			end = len(data)
		}

		// Преобразуем [][784]float32 в [][]float32
		var batch [][]float32
		for _, row := range data[i:end] {
			batch = append(batch, row[:])
		}

		batches = append(batches, batch)
	}

	return batches
}

// batchArrayLabels - функция для нарезки массива []int на батчи
func batchArrayLabels(data []int, batchSize int) [][]int {
	if batchSize <= 0 {
		panic("batchSize должен быть больше 0")
	}

	var batches [][]int
	for i := 0; i < len(data); i += batchSize {
		// Определяем конец текущего батча
		end := i + batchSize
		if end > len(data) {
			end = len(data)
		}

		// Добавляем текущий батч в результат
		batches = append(batches, data[i:end])
	}

	return batches
}

// Flatten - aункция для преобразования [][]float32 в []float32
func Flatten(batch [][]float32) []float32 {
	var result []float32
	for _, row := range batch {
		result = append(result, row...)
	}
	return result
}

// Flatten3D - функция для преобразования [][][]float32 в []float32
func Flatten3D(input [][][]float32) []float32 {
	var result []float32
	for _, matrix := range input {
		for _, row := range matrix {
			result = append(result, row...)
		}
	}
	return result
}

// ConvertToBatchedArray - превращает плоский массив в двуменрный
func ConvertToBatchedArray(flatArray []float32, batchSize int64) ([][]float32, error) {

	batchSizeInt := int(batchSize)
	
	// Проверяем, что длина массива делится на batchSize
	if len(flatArray) % batchSizeInt != 0 {
		return nil, fmt.Errorf("длина массива %d не делится на batchSize %d", len(flatArray), batchSizeInt)
	}

	// Размер каждого подмассива
	subArraySize := len(flatArray) / batchSizeInt
	batchedArray := make([][]float32, batchSize)

	for i := 0; i < batchSizeInt; i++ {
		start := i * subArraySize
		end := start + subArraySize
		batchedArray[i] = flatArray[start:end]
	}

	return batchedArray, nil
}
