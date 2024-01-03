from fastapi import FastAPI, Request

from app.routes.function import router as function_router
from app.routes.user import router as user_router

app = FastAPI()
from sqlalchemy.orm import Session
from app.database import engine, Base
from app.crud import create_function


@app.middleware("http")
async def custom_header(request: Request, call_next):
    response = await call_next(request)
    response.headers["X-Custom-Header"] = "Custom header value"
    return response


@app.on_event("startup")
async def startup_event():
    Base.metadata.create_all(bind=engine)
    with Session(engine) as db:
        # 初始化 functions
        create_function(db)


# @app.get("/")
# async def root():
#     return {"message": "hello function call."}


# app.include_router(user_router, prefix="/users", tags=["users"])
app.include_router(function_router, prefix="/functions", tags=["functions"])
