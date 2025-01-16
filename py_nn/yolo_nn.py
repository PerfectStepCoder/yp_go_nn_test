
class YOLOFashionMNIST:


    def save_to_onnx(self, file_path: str, input_size=(1, 1, 28, 28)):
        ...

    def load_onnx_model(self, model_path: str):
        ...
    

from ultralytics import YOLO

# Load the YOLO11 model
#model = YOLO("yolo11n.pt")

# Export the model to ONNX format
# path_my_model = model.export(format="onnx")
# print(f"Path my model: {path_my_model}")

# Load the exported ONNX model
onnx_model = YOLO("yolo11n.onnx")

# Run inference
results = onnx_model("https://ultralytics.com/images/bus.jpg")
print(results)
