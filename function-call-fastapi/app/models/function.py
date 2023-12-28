from sqlalchemy import Column, Integer, String, JSON

from app.database import Base


class Function(Base):
    __tablename__ = 'functions'
    id = Column(Integer, primary_key=True, index=True)
    name = Column(String(256), nullable=True, comment='名称')
    description = Column(String(256), nullable=True, comment='描述')
    parameters = Column(JSON, nullable=True, comment='参数')
