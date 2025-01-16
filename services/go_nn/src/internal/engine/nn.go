package engine

import (
	ort "github.com/yalue/onnxruntime_go"
)

// OnnxNeuralNetwork - нейронная сеть классификации изображений
type OnnxNeuralNetwork struct {
	PathToOnnxModelFile string
	PathLibOnnxRuntime  string
	InputLayer          NeuralLayer
	OutputLayer         NeuralLayer
}

type NeuralLayer struct {
	Name  string
	Shape []int64
}

func NewOnnxNeuralNetwork(pathToOnnxModelFile string, pathLibOnnxRuntime string,
	inputLayer NeuralLayer, outputLayer NeuralLayer) *OnnxNeuralNetwork {

	ort.SetSharedLibraryPath(pathLibOnnxRuntime)

	err := ort.InitializeEnvironment()
	if err != nil {
		panic(err)
	}

	return &OnnxNeuralNetwork{
		PathToOnnxModelFile: pathToOnnxModelFile,
		PathLibOnnxRuntime:  pathLibOnnxRuntime,
		InputLayer:          inputLayer,
		OutputLayer:         outputLayer,
	}
}

func (a *OnnxNeuralNetwork) Close() {
	ort.DestroyEnvironment()
}

// DetectRaw - детекция классов labels у батча с картинками images
func (a *OnnxNeuralNetwork) DetectRaw(images [][]float32) (labels [][]float32, nnError error) {
	inputData := Flatten(images)
	batchSize := int64(len(images))

	// Input
	inputSize := append([]int64{batchSize}, a.InputLayer.Shape...)
	inputShape := ort.NewShape(inputSize...)
	inputTensor, _ := ort.NewTensor(inputShape, inputData)
	defer inputTensor.Destroy()

	// Output
	outputSize := append([]int64{batchSize}, a.OutputLayer.Shape...)
	outputShape := ort.NewShape(outputSize...)
	outputTensor, _ := ort.NewEmptyTensor[float32](outputShape)
	defer outputTensor.Destroy()

	session, _ := ort.NewAdvancedSession(a.PathToOnnxModelFile,
		[]string{a.InputLayer.Name}, []string{a.OutputLayer.Name},
		[]ort.Value{inputTensor}, []ort.Value{outputTensor}, nil)
	defer session.Destroy()

	err := session.Run()
	if err != nil {
		return nil, err
	}

	outputData := outputTensor.GetData()
	result, err := ConvertToBatchedArray(outputData, batchSize)

	return result, err
}

// DetectCode - обработка батчи с картинками [N, Image] Image - одномерный массив с float картинки
func (a *OnnxNeuralNetwork) DetectCode(images [][]float32) (labelClassCode []int, nnError error) {

	if labels, err := a.DetectRaw(images); err != nil {
		return nil, err
	} else {
		labelClassCode = FindMaxIndices(labels)
	}

	return labelClassCode, nil
}

func (a *OnnxNeuralNetwork) Detect(images [][]float32) (labelClasses []string, nnError error) {

	var labelClassCode []int

	if labels, err := a.DetectRaw(images); err != nil {
		return nil, err
	} else {
		labelClassCode = FindMaxIndices(labels)
	}

	labelClasses = ConvertToClassNames(labelClassCode)

	return labelClasses, nil
}
