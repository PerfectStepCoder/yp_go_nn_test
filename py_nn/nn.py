import torch
import torch.nn as nn
import numpy as np
import onnxruntime as ort

# YOLO-подобная сеть
class YOLOFashionMNIST(nn.Module):

    def __init__(self, model_path: str = None, num_classes=10):
        """
        Конструктор нейронной сети
        """
        super(YOLOFashionMNIST, self).__init__()
        self.features = nn.Sequential(
            nn.Conv2d(1, 16, kernel_size=3, stride=1, padding=1),
            nn.ReLU(),
            nn.MaxPool2d(kernel_size=2, stride=2),
            nn.Conv2d(16, 32, kernel_size=3, stride=1, padding=1),
            nn.ReLU(),
            nn.MaxPool2d(kernel_size=2, stride=2),
        )
        self.classifier = nn.Sequential(
            nn.Linear(32 * 7 * 7, 128),
            nn.ReLU(),
            nn.Linear(128, num_classes)
        )

        self.onnx_session = None
        if model_path:
            self.load_onnx_model(model_path)

    def forward(self, x):
        """
        Прямое распространение сигнала
        """
        x = self.features(x)
        x = x.view(x.size(0), -1)
        x = self.classifier(x)
        return x
    
    def save_to_onnx(self, file_path: str, input_size=(1, 1, 28, 28)):
        """
        Сохраняет модель в формате ONNX.
        :param file_path: Путь к файлу для сохранения.
        :param input_size: Размер входного тензора (например, (1, 1, 28, 28) для FashionMNIST).
        """
        dummy_input = torch.randn(*input_size, device=next(self.parameters()).device)
        torch.onnx.export(
            self,  # Модель
            dummy_input,  # Входные данные
            file_path,  # Путь для сохранения
            export_params=True,  # Сохранить параметры обученной модели
            opset_version=11,  # Версия ONNX (можно изменить при необходимости)
            do_constant_folding=True,  # Оптимизация для постоянных выражений
            input_names=["input"],  # Имя входного слоя
            output_names=["output"],  # Имя выходного слоя
            dynamic_axes={"input": {0: "batch_size"}, "output": {0: "batch_size"}},  # Поддержка переменного размера батча
        )
        print(f"Model successfully saved to {file_path}") 

    def load_onnx_model(self, model_path: str):
        """
        Метод для загрузки модели в формате ONNX.
        :param model_path: Путь к файлу модели в формате ONNX.
        """
        try:
            self.onnx_session = ort.InferenceSession(model_path)
            print(f"ONNX модель успешно загружена из {model_path}")
        except Exception as e:
            print(f"Ошибка загрузки ONNX модели: {e}")
            raise
    
    def predict_onnx(self, input_data: torch.Tensor) -> np.ndarray:
        """
        Метод для выполнения предсказания с использованием ONNX модели.
        :param input_data: Входные данные в формате Tensor.
        :return: Предсказания модели в формате numpy.
        """
        if self.onnx_session is None:
            raise ValueError("ONNX модель не загружена. Сначала вызовите load_onnx_model().")

        # Преобразование тензора в numpy
        input_data_np = input_data.numpy().astype(np.float32)

        # Получение имени входного и выходного тензоров
        input_name = self.onnx_session.get_inputs()[0].name
        output_name = self.onnx_session.get_outputs()[0].name

        # Выполнение предсказания
        predictions = self.onnx_session.run([output_name], {input_name: input_data_np})[0]
        return predictions
