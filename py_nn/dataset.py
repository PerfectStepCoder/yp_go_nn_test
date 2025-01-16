import pandas as pd
from torch.utils.data import DataLoader
from torchvision import datasets, transforms


def get_dataloaders(batch_size=32) -> tuple[DataLoader, DataLoader]:

    # Подготовка данных
    transform = transforms.Compose([transforms.ToTensor(), transforms.Normalize((0.5,), (0.5,))])

    train_dataset = datasets.FashionMNIST(root="./data", train=True, transform=transform, download=True)
    test_dataset = datasets.FashionMNIST(root="./data", train=False, transform=transform, download=True)

    train_loader = DataLoader(train_dataset, batch_size=batch_size, shuffle=True)
    test_loader = DataLoader(test_dataset, batch_size=batch_size, shuffle=False)

    return train_loader, test_loader


def export_test(file_path: str = 'fashion_mnist_test.csv'):
    
    # Трансформация данных
    transform = transforms.Compose([transforms.ToTensor(), transforms.Normalize((0.5,), (0.5,))])

    # Загружаем тестовый датасет
    test_dataset = datasets.FashionMNIST(root='./data', train=False, transform=transform, download=True)

    # Преобразование изображений и меток в DataFrame
    data = []
    for image, label in test_dataset:
        flattened_image = image.view(-1).tolist()  # Преобразование 28x28 в 1x784
        data.append([label] + flattened_image)

    columns = ['label'] + [f'pixel_{i}' for i in range(28*28)]
    df = pd.DataFrame(data, columns=columns)

    # Сохраняем в CSV
    df.to_csv(file_path, index=False)
    print(f"Датасет успешно экспортирован в {file_path}")