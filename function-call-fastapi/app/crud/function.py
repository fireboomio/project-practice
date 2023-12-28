import json

from sqlalchemy.orm import Session

from app.models.function import Function
from app.utils.auto_tool_generator import AutoFunctionGenerator
from app.utils.functions import (
    gpt_dall_e,
    google_search
)


def get_functions(db: Session):
    db_functions = db.query(
        Function.name,
        Function.description,
        Function.parameters
    ).all()

    functions = []
    for func in db_functions:
        function_data = {
            "name": func.name,
            "description": func.description,
            "parameters": func.parameters
        }
        functions.append(function_data)

    return functions


def create_function(db: Session):
    functions_list = [gpt_dall_e, google_search]
    generator = AutoFunctionGenerator(functions_list)
    function_descriptions = generator.auto_generate()
    for function in function_descriptions:
        db_function = Function(
            name=function.get('name'),
            description=function.get('description'),
            parameters=function.get('parameters'),
        )
        db.add(db_function)
        db.commit()
        db.refresh(db_function)
    functions = get_functions(db)
    return functions


def execute_function(name, arguments):
    print('方法：', name)
    print('参数：', arguments)

    # 检查函数是否存在
    if name in globals():
        # 获取函数引用
        func = globals()[name]

        # 将参数字符串转换为 JSON 对象
        try:
            arguments = json.loads(arguments)
        except json.JSONDecodeError:
            return {"error": "Invalid JSON format in arguments"}

        # 调用函数
        try:
            response = func(**arguments)
            result = {"result": response}
            return result
        except Exception as e:
            return {"error": str(e)}
    else:
        return {"error": f"Function '{name}' not found"}
