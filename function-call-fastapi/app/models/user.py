from sqlalchemy import Boolean, Column, Integer, String

from app.database import Base


class User(Base):
    __tablename__ = "users"
    id = Column(Integer, primary_key=True, index=True)
    email = Column(String(200), unique=True, index=True)
    full_name = Column(String(200), index=True)
    hashed_password = Column(String(200))
    is_active = Column(Boolean, default=True)
