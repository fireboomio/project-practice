from typing import Optional

from pydantic import BaseModel


class UserBase(BaseModel):
    email: str
    full_name: Optional[str] = None
    password: str


class UserCreate(UserBase):
    pass


class UserUpdate(UserBase):
    password: Optional[str] = None


class User(UserBase):
    id: int
    is_active: bool

    class Config:
        orm_mode = True
