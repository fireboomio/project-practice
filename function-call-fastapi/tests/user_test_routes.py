from fastapi.testclient import TestClient
from app.main import app
from app.database import SessionLocal
from app.schemas.user import User, UserCreate
from app.crud.user import create_user, get_user_by_email

client = TestClient(app)

# Testing POST / route
def test_create_user():
    # Limpiar la base de datos
    db = SessionLocal()
    db.query(User).delete()
    db.commit()

    # Enviar solicitud POST con datos de usuario válidos
    user = {"email": "test@example.com", "password": "password"}
    response = client.post("/", json=user)

    # Verificar que la respuesta sea exitosa
    assert response.status_code == 201
    assert response.json()["email"] == "test@example.com"

    # Verificar que el usuario se haya creado en la base de datos
    db_user = get_user_by_email(db, email="test@example.com")
    assert db_user is not None

    # Limpiar la base de datos
    db.query(User).delete()
    db.commit()

# Testing GET /{user_id} route
def test_read_user():
    # Limpiar la base de datos y crear un usuario de prueba
    db = SessionLocal()
    db.query(User).delete()
    db_user = create_user(db, user=UserCreate(email="test@example.com", password="password"))
    db.commit()

    # Enviar solicitud GET para el usuario de prueba
    response = client.get(f"/{db_user.id}")

    # Verificar que la respuesta sea exitosa y contenga los datos del usuario
    assert response.status_code == 200
    assert response.json()["email"] == "test@example.com"

    # Limpiar la base de datos
    db.query(User).delete()
    db.commit()