from fastapi.testclient import TestClient

from app.crud.user import create_user, get_user_by_email
from app.database import SessionLocal
from app.main import app
from app.schemas.user import User, UserCreate

client = TestClient(app)


def test_create_user():
    db = SessionLocal()
    db.query(User).delete()
    db.commit()

    user = {"email": "test@example.com", "password": "password"}
    response = client.post("/", json=user)

    assert response.status_code == 201
    assert response.json()["email"] == "test@example.com"

    db_user = get_user_by_email(db, email="test@example.com")
    assert db_user is not None

    db.query(User).delete()
    db.commit()


def test_read_user():
    db = SessionLocal()
    db.query(User).delete()
    db_user = create_user(db, user=UserCreate(email="test@example.com", password="password"))
    db.commit()

    response = client.get(f"/{db_user.id}")

    assert response.status_code == 200
    assert response.json()["email"] == "test@example.com"

    db.query(User).delete()
    db.commit()
