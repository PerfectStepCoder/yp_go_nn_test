import onnx

# Загрузка модели
model = onnx.load("/Users/dmitrii/EducationProjects/yp_go_nn/services/go_nn/src/cmd/yolo_fashion_mnist.onnx")
print("Входные слои:")
for input in model.graph.input:
    print(input.name)

print("Выходные слои:")
for output in model.graph.output:
    print(output.name)
