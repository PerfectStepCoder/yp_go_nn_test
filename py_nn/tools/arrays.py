import numpy as np

def bytes_to_float_matrix(data: bytes, rows: int, cols: int) -> np.ndarray:
    """
    Преобразует массив bytes в двумерный массив float.

    :param data: Массив байт, содержащий float32.
    :param rows: Количество строк в результирующем массиве.
    :param cols: Количество столбцов в результирующем массиве.
    :return: Двумерный массив numpy.ndarray с float32.
    """
    # Проверяем, что размер данных соответствует ожидаемому размеру
    expected_size = rows * cols * 4  # float32 занимает 4 байта
    if len(data) != expected_size:
        raise ValueError(f"Размер данных {len(data)} байт не соответствует ожидаемому {expected_size} байтам.")
    
    # Преобразуем байты в массив float32
    flat_array = np.frombuffer(data, dtype=np.float32)
    
    # Преобразуем в двумерный массив
    return flat_array.reshape((rows, cols))


def find_max_indices(array: np.ndarray) -> list[int]:
    """
    Находит индексы максимальных элементов в каждой строке массива.
    
    :param array: Двумерный массив numpy.ndarray
    :return: Список индексов максимальных элементов для каждой строки
    """
    if not isinstance(array, np.ndarray):
        raise ValueError("Input должен быть numpy.ndarray")
    if len(array.shape) != 2:
        raise ValueError("Input должен быть двумерным массивом")

    # Находим индексы максимальных элементов в каждой строке
    max_indices = np.argmax(array, axis=1).astype(np.int32)

    return max_indices.tolist()