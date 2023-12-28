from fastapi import APIRouter, Depends, HTTPException, status
from fastapi.responses import JSONResponse
from sqlalchemy.orm import Session

from app.crud import get_user_by_email, create_user, get_user
from app.database import get_db
from app.schemas import UserCreate

router = APIRouter()


@router.post("/create", status_code=status.HTTP_201_CREATED)
def create_user(user: UserCreate, db: Session = Depends(get_db)):
    db_user = get_user_by_email(db, email=user.email)
    if db_user:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail="Email already registered",
        )
    return create_user(db=db, user=user)


@router.get("/{user_id}")
async def read_user(user_id: int, db: Session = Depends(get_db)):
    try:
        db_user = get_user(db, user_id=user_id)
        print(db_user)
        if db_user is None:
            return JSONResponse(content={"message": "User not found"},
                                status_code=status.HTTP_404_NOT_FOUND)
        return db_user
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail="Error: " + str(e),
        )
