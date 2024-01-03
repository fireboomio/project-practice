import inspect
import json

from app.utils.gpt import GPT


class AutoFunctionGenerator:
    """
    AutoFunctionGenerator 类用于自动生成一系列功能函数的 JSON Schema 描述。
    该类通过调用 OpenAI API，采用 Few-shot learning 的方式来生成这些描述。

    属性:
    - functions_list (list): 一个包含多个功能函数的列表。
    - max_attempts (int): 最大尝试次数，用于处理 API 调用失败的情况。

    方法:
    - __init__ : 初始化 AutoFunctionGenerator 类。
    - generate_function_descriptions : 自动生成功能函数的 JSON Schema 描述。
    - _call_openai_api : 调用 OpenAI API。
    - auto_generate : 自动生成功能函数的 JSON Schema 描述，并处理任何异常。
    """

    def __init__(self, functions_list, max_attempts=3):
        """
        初始化 AutoFunctionGenerator 类。

        参数:
        - functions_list (list): 一个包含多个功能函数的列表。
        - max_attempts (int): 最大尝试次数。
        """
        self.functions_list = functions_list
        self.max_attempts = max_attempts

    def generate_function_descriptions(self):
        """
        自动生成功能函数的 JSON Schema 描述。

        返回:
        - list: 包含 JSON Schema 描述的列表。
        """
        # 创建空列表，保存每个功能函数的JSON Schema描述
        functions = []

        for function in self.functions_list:
            # 读取指定函数的函数说明
            function_description = inspect.getdoc(function)

            # 读取函数的函数名
            function_name = function.__name__

            # 定义system role的Few-shot提示
            system_Q = "你是一位优秀的数据分析师，现在有一个函数的详细声明如下：%s" % function_description
            system_A = "计算年龄总和的函数，该函数从一个特定格式的JSON字符串中解析出DataFrame，然后计算所有人的年龄总和并以JSON格式返回结果。\
                        \n:param input_json: 必要参数，要求字符串类型，表示含有个体年龄数据的JSON格式字符串 \
                        \n:return: 计算完成后的所有人年龄总和，返回结果为JSON字符串类型对象"

            # 定义user role的Few-shot提示
            user_Q = "请根据这个函数声明，为我生成一个JSON Schema对象描述。这个描述应该清晰地标明函数的输入和输出规范。具体要求如下：\
                      1. 提取函数名称：%s，并将其用作JSON Schema中的'name'字段  \
                      2. 在JSON Schema对象中，设置函数的参数类型为'object'.\
                      3. 'properties'字段如果有参数，必须表示出字段的描述. \
                      4. 从函数声明中解析出函数的描述，并在JSON Schema中以中文字符形式表示在'description'字段.\
                      5. 识别函数声明中哪些参数是必需的，然后在JSON Schema的'required'字段中列出这些参数. \
                      6. 输出的应仅为符合上述要求的JSON Schema对象内容,不需要任何上下文修饰语句. " % function_name

            user_A = "{'name': 'calculate_total_age_function', \
                               'description': '计算年龄总和的函数，从给定的JSON格式字符串（按'split'方向排列）中解析出DataFrame，计算所有人的年龄总和，并以JSON格式返回结果。 \
                               'parameters': {'type': 'object', \
                                              'properties': {'input_json': {'description': '执行计算年龄总和的数据集', 'type': 'string'}}, \
                                              'required': ['input_json']}}"

            # 定义输入
            system_message = "你是一位优秀的数据分析师，现在有一个函数的详细声明如下：%s" % function_description
            user_message = "请根据这个函数声明，为我生成一个JSON Schema对象描述。这个描述应该清晰地标明函数的输入和输出规范。具体要求如下：\
                            1. 提取函数名称：%s，并将其用作JSON Schema中的'name'字段  \
                            2. 在JSON Schema对象中，设置函数的参数类型为'object'.\
                            3. 'properties'字段如果有参数，必须表示出字段的描述. \
                            4. 从函数声明中解析出函数的描述，并在JSON Schema中以中文字符形式表示在'description'字段.\
                            5. 识别函数声明中哪些参数是必需的，然后在JSON Schema的'required'字段中列出这些参数. \
                            6. 输出的应仅为符合上述要求的JSON Schema对象内容,不需要任何上下文修饰语句. " % function_name

            messages = [
                {"role": "system", "content": "Q:" + system_Q + user_Q + "A:" + system_A + user_A},

                {"role": "user", "content": 'Q:' + system_message + user_message}
            ]

            response = self._call_openai_api(messages)
            # 获取JSON格式的字符串，并去除Markdown格式
            function_json_str = response["choices"][0]["message"]["content"]
            function_json_str = function_json_str.strip('```json\n').strip('\n```')

            # 解析JSON字符串为Python字典
            try:
                function_json = json.loads(function_json_str)
                functions.append(function_json)
            except json.JSONDecodeError as e:
                print(f"JSON parsing error: {e}")
                functions.append({})
        return functions

    def _call_openai_api(self, messages):
        """
        私有方法，用于调用 OpenAI API。

        参数:
        - messages (list): 包含 API 所需信息的消息列表。

        返回:
        - object: API 调用的响应对象。
        """
        # 请根据您的实际情况修改此处的 API 调用
        return GPT.chat(messages)

    def auto_generate(self):
        """
        自动生成功能函数的 JSON Schema 描述，并处理任何异常。

        返回:
        - list: 包含 JSON Schema 描述的列表。

        异常:
        - 如果达到最大尝试次数，将抛出异常。
        """
        attempts = 0
        while attempts < self.max_attempts:
            try:
                functions = self.generate_function_descriptions()
                return functions
            except Exception as e:
                attempts += 1
                print(f"Error occurred: {e}")
                if attempts >= self.max_attempts:
                    print("Reached maximum number of attempts. Terminating.")
                    raise
                else:
                    print("Retrying...")
