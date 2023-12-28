from fastapi import FastAPI, Request

from app.routes.function import router as function_router
from app.routes.user import router as user_router

app = FastAPI()


@app.middleware("http")
async def custom_header(request: Request, call_next):
    response = await call_next(request)
    response.headers["X-Custom-Header"] = "Custom header value"
    return response


@app.get("/")
async def root():
    return {"message": "hello function call."}


app.include_router(user_router, prefix="/users", tags=["users"])
app.include_router(function_router, prefix="/functions", tags=["functions"])
