from typing import Optional, List, Dict

from pydantic import BaseModel, Field


class FunctionExecute(BaseModel):
    name: str
    arguments: str


class PropertyModel(BaseModel):
    type: str
    description: str


class ParametersModel(BaseModel):
    type: str = Field(default="object")
    required: Optional[List[str]] = None
    properties: Dict[str, PropertyModel]


class FunctionModel(BaseModel):
    name: str
    description: Optional[str] = None
    parameters: Optional[ParametersModel] = None


class ExecutionResult(BaseModel):
    result: str
