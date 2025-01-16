import torch
import torch.nn as nn
import torch.optim as optim
import numpy as np
from dataset import get_dataloaders, export_test
from nn import YOLOFashionMNIST
from test import test
from train import train


# Настройки
batch_size = 64
epochs = 5
learning_rate = 0.001
onnx_file_path = "/Users/dmitrii/EducationProjects/yp_go_nn/services/go_nn/src/models/yolo_fashion_mnist.onnx"
device = torch.device("cuda" if torch.cuda.is_available() else "cpu")
print(f"Device: {device}")

train_loader, test_loader = get_dataloaders(batch_size=batch_size)

# Инициализация модели
model = YOLOFashionMNIST().to(device)
criterion = nn.CrossEntropyLoss()
optimizer = optim.Adam(model.parameters(), lr=learning_rate)


# Запуск обучения и тестирования
if __name__ == "__main__":

    print("0 - export test dataset\n1 - train and test model\n2 - test saved model")

    match int(input("Enter cmd:")):
        case 0:
            export_test()
        case 1:
            train(model, epochs, train_loader, optimizer, criterion, device)
            model.save_to_onnx(onnx_file_path)
            test(model, test_loader, device)
        case 2:
            model.load_onnx_model(onnx_file_path)
            # Выбор одного батча данных
            images, labels = next(iter(test_loader))
            # Выполнение предсказания
            predictions = model.predict_onnx(images)
            # Вывод результата
            predicted_label = np.argmax(predictions, axis=1)
            print(f"Предсказанная метка: {predicted_label[0]}, Истинная метка: {labels[0].item()}")
