package engine

import (
	"bytes"
	"fmt"
	"encoding/binary"
	//"errors"
	"math"
)

func FindMaxIndices(labels [][]float32) []int {
	// Массив для хранения индексов максимальных элементов
	maxIndices := make([]int, len(labels))

	// Перебираем каждый ряд массива
	for i, row := range labels {
		maxIndex := 0
		maxValue := row[0]

		// Найти индекс максимального элемента в текущем ряду
		for j, value := range row {
			if value > maxValue {
				maxValue = value
				maxIndex = j
			}
		}

		// Добавляем индекс в результат
		maxIndices[i] = maxIndex
	}

	return maxIndices
}

// Функция для сравнения двух массивов
func CompareArrays(arr1, arr2 []int) bool {
	// Проверяем, что массивы одинаковой длины
	if len(arr1) != len(arr2) {
		return false
	}

	// Сравниваем каждый элемент
	for i := range arr1 {
		if arr1[i] != arr2[i] {
			return false
		}
	}

	return true
}

// Функция для подсчета процента совпадений
func CalculateMatchPercentage(arr1, arr2 []int) float64 {
	// Проверяем, что массивы одинаковой длины
	if len(arr1) != len(arr2) {
		return 0 // Если массивы разной длины, процент совпадений невозможен
	}

	// Считаем количество совпавших элементов
	matches := 0
	for i := range arr1 {
		if arr1[i] == arr2[i] {
			matches++
		}
	}

	// Рассчитываем процент совпадений
	percentage := (float64(matches) / float64(len(arr1))) * 100
	return percentage
}

// Float32MatrixToBytes
func Float32MatrixToBytes(matrix [][]float32) ([]byte, error) {
	var buf bytes.Buffer

	// Проходим по строкам и колонкам матрицы
	for _, row := range matrix {
		for _, value := range row {
			// Конвертируем float32 в []byte и записываем в буфер
			if err := binary.Write(&buf, binary.LittleEndian, value); err != nil {
				return nil, fmt.Errorf("failed to write float32: %w", err)
			}
		}
	}

	return buf.Bytes(), nil
}

// BytesToFloat32Matrix - преобразование массива байт в исходный массив
func BytesToFloat32Matrix(data []byte, rows, cols int) ([][]float32, error) {
	// Проверяем, что длина данных соответствует размерам матрицы
	expectedLength := rows * cols * 4 // 4 байта на float32
	if len(data) != expectedLength {
		return nil, fmt.Errorf("invalid data length: expected %d, got %d", expectedLength, len(data))
	}

	// Создаём матрицу нужного размера
	matrix := make([][]float32, rows)
	for i := range matrix {
		matrix[i] = make([]float32, cols)
	}

	// Читаем данные из []byte
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			offset := (i*cols + j) * 4 // Смещение для текущего элемента
			uintValue := binary.LittleEndian.Uint32(data[offset : offset+4])  // Читаем uint32
			matrix[i][j] = math.Float32frombits(uintValue)                    // Преобразуем uint32 в float32
		}
	}

	return matrix, nil
}

// GetMatrixSize вычисляет количество строк и столбцов в массиве [][]float32
func GetMatrixSize(matrix [][]float32) (rows, cols int) {
	rows = len(matrix) // Количество строк — это длина внешнего массива
	if rows > 0 {
		cols = len(matrix[0]) // Количество столбцов — это длина первой строки
	}
	return rows, cols
}