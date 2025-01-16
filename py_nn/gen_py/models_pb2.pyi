from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class TaskOneRequest(_message.Message):
    __slots__ = ("taskUID", "image", "height", "width")
    TASKUID_FIELD_NUMBER: _ClassVar[int]
    IMAGE_FIELD_NUMBER: _ClassVar[int]
    HEIGHT_FIELD_NUMBER: _ClassVar[int]
    WIDTH_FIELD_NUMBER: _ClassVar[int]
    taskUID: str
    image: bytes
    height: int
    width: int
    def __init__(self, taskUID: _Optional[str] = ..., image: _Optional[bytes] = ..., height: _Optional[int] = ..., width: _Optional[int] = ...) -> None: ...

class TaskOneResponse(_message.Message):
    __slots__ = ("taskUID", "className")
    TASKUID_FIELD_NUMBER: _ClassVar[int]
    CLASSNAME_FIELD_NUMBER: _ClassVar[int]
    taskUID: str
    className: str
    def __init__(self, taskUID: _Optional[str] = ..., className: _Optional[str] = ...) -> None: ...

class TaskBatchRequest(_message.Message):
    __slots__ = ("taskUID", "images", "height", "width")
    TASKUID_FIELD_NUMBER: _ClassVar[int]
    IMAGES_FIELD_NUMBER: _ClassVar[int]
    HEIGHT_FIELD_NUMBER: _ClassVar[int]
    WIDTH_FIELD_NUMBER: _ClassVar[int]
    taskUID: str
    images: _containers.RepeatedScalarFieldContainer[bytes]
    height: int
    width: int
    def __init__(self, taskUID: _Optional[str] = ..., images: _Optional[_Iterable[bytes]] = ..., height: _Optional[int] = ..., width: _Optional[int] = ...) -> None: ...

class TaskBatchResponse(_message.Message):
    __slots__ = ("taskUID", "classNames")
    TASKUID_FIELD_NUMBER: _ClassVar[int]
    CLASSNAMES_FIELD_NUMBER: _ClassVar[int]
    taskUID: str
    classNames: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, taskUID: _Optional[str] = ..., classNames: _Optional[_Iterable[str]] = ...) -> None: ...

class TaskBatchCodeResponse(_message.Message):
    __slots__ = ("taskUID", "classCodes")
    TASKUID_FIELD_NUMBER: _ClassVar[int]
    CLASSCODES_FIELD_NUMBER: _ClassVar[int]
    taskUID: str
    classCodes: _containers.RepeatedScalarFieldContainer[int]
    def __init__(self, taskUID: _Optional[str] = ..., classCodes: _Optional[_Iterable[int]] = ...) -> None: ...

class ServiceInfoNN(_message.Message):
    __slots__ = ("name", "description", "version")
    NAME_FIELD_NUMBER: _ClassVar[int]
    DESCRIPTION_FIELD_NUMBER: _ClassVar[int]
    VERSION_FIELD_NUMBER: _ClassVar[int]
    name: str
    description: str
    version: str
    def __init__(self, name: _Optional[str] = ..., description: _Optional[str] = ..., version: _Optional[str] = ...) -> None: ...
