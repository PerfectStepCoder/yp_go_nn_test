package engine

import (
	"bytes"
	"encoding/binary"
	"fmt"

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

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// Функция для сравнения двух массивов
func CompareArrays[T Integer](arr1, arr2 []T) bool {
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
func CalculateMatchPercentage[T Integer](arr1, arr2 []T) float64 {
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
			offset := (i*cols + j) * 4                                       // Смещение для текущего элемента
			uintValue := binary.LittleEndian.Uint32(data[offset : offset+4]) // Читаем uint32
			matrix[i][j] = math.Float32frombits(uintValue)                   // Преобразуем uint32 в float32
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

// GetFirstImage из массива извлекает только первое изображение
func GetFirstImage(input [][]float32) [][]float32 {

	// Проверяем, что input содержит хотя бы одну строку
	if len(input) == 0 {
		return [][]float32{} // Возвращаем пустой массив, если input пуст
	}

	// Вычисляем ширину (количество столбцов)
	width := len(input[0])

	// Создаем массив [][]float32 с динамической длиной строки
	var output [][]float32 = make([][]float32, 1)
	output[0] = make([]float32, width)

	// Копируем данные из первой строки input в output
	for i := 0; i < width; i++ {
		output[0][i] = input[0][i]
	}

	return output
}

// BytesToFloat32Slice - массив []bytes -> []float32
func BytesToFloat32Slice(data []byte) ([]float32, error) {
	// Длина массива float32
	if len(data)%4 != 0 {
		return nil, fmt.Errorf("invalid byte slice length: %d (must be multiple of 4)", len(data))
	}

	// Вычисляем длину результирующего массива
	numFloats := len(data) / 4
	result := make([]float32, numFloats)

	// Чтение данных
	buf := bytes.NewReader(data)
	for i := 0; i < numFloats; i++ {
		var f float32
		err := binary.Read(buf, binary.LittleEndian, &f)
		if err != nil {
			return nil, err
		}
		result[i] = f
	}
	return result, nil
}

func IntToInt32Slice(input []int) []int32 {
	// Создаем массив []int32 с той же длиной
	result := make([]int32, len(input))

	// Копируем элементы с явным преобразованием
	for i, v := range input {
		result[i] = int32(v)
	}

	return result
}

func Float32ToBytes2D(input [][]float32) ([][]byte, error) {
	// Создаём выходной массив [][]byte с такой же структурой
	output := make([][]byte, len(input))

	for i, row := range input {
		// Создаём срез []byte для каждой строки
		rowBytes := new(bytes.Buffer)
		for _, value := range row {
			// Записываем float32 в байты
			err := binary.Write(rowBytes, binary.LittleEndian, value)
			if err != nil {
				return nil, err
			}
		}
		output[i] = rowBytes.Bytes()
	}

	return output, nil
}
