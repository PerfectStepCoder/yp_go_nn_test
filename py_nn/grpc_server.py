from concurrent import futures
import grpc
import torch
import sys
import os
import logging
import struct
import numpy as np

# Добавляем родительский каталог в sys.path
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), '.')))
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), './gen_py')))
print(sys.path)

from tools.arrays import bytes_to_float_matrix
from tools.mnist import fashion_mnist_classes
from nn import YOLOFashionMNIST

from gen_py import models_pb2
from gen_py import server_pb2_grpc


class NNService(server_pb2_grpc.ClassifyNNServicer):
    
    def __init__(self, model_path):
        super().__init__()
        self.model_path = model_path
        device = torch.device("cuda" if torch.cuda.is_available() else "cpu")
        self._nn = YOLOFashionMNIST().to(device)
        self._nn.load_onnx_model(self.model_path)

    def GetInfo(self, request, context):
        response = models_pb2.ServiceInfoNN(
            name="YoloNN", description="[Python] Trained Fashion MNIST", version="1.0.0.0"   
        )
        return response
    
    def CreateBatchCodeTask(self, request: models_pb2.TaskBatchRequest, context):
        response = models_pb2.TaskBatchCodeResponse(
            taskUID=request.taskUID
        )
        batch_size = len(request.images)
        input_batch = []
        for img in request.images:
            input_batch.append(bytes_to_float_matrix(img, request.height, request.width))
        tensor = torch.tensor(input_batch)
        tensor = tensor.view(batch_size, 1, request.height, request.width)
        
        for class_code in self._nn.detect_code_onnx(tensor):
            response.classCodes.append(class_code)
        
        return response
    
    def CreateBatchTask(self, request: models_pb2.TaskBatchRequest, context):
        
        response = models_pb2.TaskBatchResponse(
            taskUID=request.taskUID
        )
        
        batch_size = len(request.images)
        input_batch = bytes_to_float_matrix(request.images, batch_size, request.width)
        tensor = torch.tensor(input_batch)
        tensor = tensor.view(batch_size, 1, request.height, request.width)
        class_codes = self._nn.detect_code_onnx(input)
        response.classNames = fashion_mnist_classes(class_codes)
        return response
    
    def serve(model_path):
        server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
        server_pb2_grpc.add_ClassifyNNServicer_to_server(NNService(model_path), server)
        server.add_insecure_port('localhost:50051')
        print("Server started on port 50051")
        server.start()
        server.wait_for_termination()


def main():
    
    model_path="/Users/dmitrii/EducationProjects/yp_go_nn/services/go_nn/src/models/yolo_fashion_mnist.onnx"
    NNService.serve(model_path)

if __name__ == "__main__":

    main()
