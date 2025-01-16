package engine

// Тип для классов
type FashionClass int

// Константы для классов
const (
	TShirtTop FashionClass = iota // 0
	Trouser                       // 1
	Pullover                      // 2
	Dress                         // 3
	Coat                          // 4
	Sandal                        // 5
	Shirt                         // 6
	Sneaker                       // 7
	Bag                           // 8
	AnkleBoot                     // 9
)

// Названия классов в строковом виде
var fashionClassNames = map[FashionClass]string{
	TShirtTop: "T-shirt/Top",
	Trouser:   "Trouser",
	Pullover:  "Pullover",
	Dress:     "Dress",
	Coat:      "Coat",
	Sandal:    "Sandal",
	Shirt:     "Shirt",
	Sneaker:   "Sneaker",
	Bag:       "Bag",
	AnkleBoot: "Ankle Boot",
}

// Метод для преобразования FashionClass в строку
func (fc FashionClass) String() string {
	return fashionClassNames[fc]
}

// Преобразование массива целых чисел в массив строк
func ConvertToClassNames(indices []int) []string {
	classNames := make([]string, len(indices))
	for i, index := range indices {
		// Преобразуем индекс в тип FashionClass для извлечения названия
		classNames[i] = fashionClassNames[FashionClass(index)]
	}
	return classNames
}
