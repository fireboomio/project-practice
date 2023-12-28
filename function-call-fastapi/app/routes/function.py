from typing import List

from fastapi import APIRouter, Depends, HTTPException, status
from fastapi.responses import JSONResponse
from sqlalchemy.orm import Session

from app.crud import get_functions, create_function, execute_function
from app.database import get_db
from app.schemas import FunctionExecute, FunctionModel, ExecutionResult

router = APIRouter()


@router.post("/create", status_code=status.HTTP_201_CREATED)
async def create(db: Session = Depends(get_db)):
    functions = create_function(db)
    return functions


@router.get("/all", response_model=List[FunctionModel])
async def read(db: Session = Depends(get_db)):
    try:
        functions = get_functions(db)
        print(functions)
        if functions is None:
            return JSONResponse(content={"message": "not found"}, status_code=status.HTTP_404_NOT_FOUND)
        return functions
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail="Error: " + str(e),
        )


@router.post("/execute", response_model=ExecutionResult, status_code=status.HTTP_201_CREATED)
async def execute(function: FunctionExecute):
    result = execute_function(function.name, function.arguments)
    return result
