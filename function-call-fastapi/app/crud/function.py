import inspect
import json

from sqlalchemy.orm import Session

from app.models.function import Function
from app.utils import functions as utils_functions
from app.utils.auto_tool_generator import AutoFunctionGenerator


def get_functions(db: Session):
    db_functions = db.query(Function.name, Function.description, Function.parameters).all()
    return [
        {"name": func.name, "description": func.description, "parameters": func.parameters}
        for func in db_functions
    ]


def get_function_by_name(db: Session, name: str):
    return db.query(Function).filter(Function.name == name).first()


def create_function(db: Session):
    all_functions = inspect.getmembers(utils_functions, inspect.isfunction)

    for name, func in all_functions:
        if not get_function_by_name(db, name):
            generator = AutoFunctionGenerator([func])
            function_descriptions = generator.auto_generate()
            for function in function_descriptions:
                db_function = Function(
                    name=function.get('name'),
                    description=function.get('description'),
                    parameters=function.get('parameters'),
                )
                db.add(db_function)
            db.commit()


def execute_function(name, arguments):
    print('方法：', name)
    print('参数：', arguments)

    function_map = dict(inspect.getmembers(utils_functions, inspect.isfunction))

    func = function_map.get(name)
    if not func:
        return {"error": f"Function '{name}' not found"}

    try:
        arguments = json.loads(arguments)
        response = func(**arguments)
        return {"result": response}
    except json.JSONDecodeError:
        return {"error": "Invalid JSON format in arguments"}
    except Exception as e:
        return {"error": str(e)}
