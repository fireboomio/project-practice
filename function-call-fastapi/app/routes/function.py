import asyncio
from typing import List

from fastapi import APIRouter, Depends, HTTPException, status
from fastapi.responses import JSONResponse
from sqlalchemy.orm import Session

from app.crud import get_functions, execute_function
from app.database import get_db
from app.schemas import FunctionExecute, FunctionModel, ExecutionResult

router = APIRouter()


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
    loop = asyncio.get_event_loop()
    result = await loop.run_in_executor(None, execute_function, function.name, function.arguments)
    print("返回：", result)
    return result
