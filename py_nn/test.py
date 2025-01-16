import torch

# Тестирование модели
def test(model, test_loader, device, verbose=True):

    model.eval()

    correct = 0
    total = 0
    with torch.no_grad():
        for images, labels in test_loader:
            #print(images)
            images, labels = images.to(device), labels.to(device)
            outputs = model(images)
            print(outputs)
            _, predicted = torch.max(outputs.data, 1)
            total += labels.size(0)
            correct += (predicted == labels).sum().item()

    if verbose:
        print(f"Accuracy on test data: {100 * correct / total:.2f}%")