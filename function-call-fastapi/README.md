<p align="center">
  <a href=""><img src="https://raw.githubusercontent.com/songquanpeng/one-api/main/web/public/logo.png" width="150" height="150" alt="one-api logo"></a>
</p>

<div align="center">

# Function Call

_✨ 为 GPT functioncalling 提供 functions ✨_

</div>

## 部署

### 基于 Docker Compose 进行部署（确保已安装 Docker 和 Docker Compose）

1. 先配置服务所需的环境变量，根据 env.example 文件创建 .env 文件，并填写相应的环境变量

```shell
# 配置 mysql
MYSQL_ROOT_PASSWORD=root用户密码
MYSQL_PASSWORD=用户密码
MYSQL_USER=用户名
MYSQL_DATABASE=数据库名称
MYSQL_PORT=数据库端口
MYSQL_HOST=数据库地址
# 配置 SQLAlchemy 链接，以 mysql 举例
DB_URL=mysql+pymysql://root:${MYSQL_ROOT_PASSWORD}@${MYSQL_HOST}:${MYSQL_PORT}/${MYSQL_DATABASE}

# 配置 oneapi
ONE_API_KEY=
ONE_API_URL=

# 配置 google serper 
GOOGLE_API_KEY=
GOOGLE_API_URL=
```

2. 启动服务

```shell
# 构建和启动容器
docker-compose up -d

# 查看容器启动状态
docker-compose ps

# 查看 app 服务日志
docker-compose logs -f app
```

### 手动部署

```shell
1. 配置 python 环境
# 安装 virtualenv
pip install virtualenv

# 使用 virtualenv 创建虚拟环境
virtualenv venv

# 激活虚拟环境
source venv/bin/activate

# 安装依赖
pip install --no-cache-dir -r requirements.txt

# 安装数据库迁移工具
pip install alembic

# 迁移数据库到最新
alembic upgrade head
```

2. 启动服务

```shell
cuvicorn app.main:app --host 0.0.0.0 --port 5001
```

## 访问

访问 [http://localhost:5001/docs](http://localhost:5001/docs)，查看 Swagger 文档。

## 使用

架构图

![Alt text](https://ft-dev.oss-cn-shanghai.aliyuncs.com/WechatIMG150.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=LTAI5tP8iEBMfrzVNSW7k148%2F20231228%2Foss-cn-shanghai%2Fs3%2Faws4_request&X-Amz-Date=20231228T080054Z&X-Amz-Expires=86400&X-Amz-SignedHeaders=host&response-content-disposition=attachment%3Bfilename%3D%22WechatIMG150.jpg%22&X-Amz-Signature=5424c5be7af147154c2bda230765d37eabadf6d66e41358a6c8b01b3709adfaf)

系统本身开箱即用。

你只需要在 `utils/functions.py` 下定义你所需要的函数，确保函数返回一个 `字符串类型` 的结果，最关键是完善函数的描述：

包括：函数功能描述、参数（类型和描述）、返回（类型和描述）。

系统内置函数：

```python
# GPT 画图
def gpt_dall_e(prompt, model='dall-e-3', n=1, size='1024x1024'):
    """
    根据给定的提示（prompt），使用指定的 DALL-E 模型生成图像，并返回图像的 URL。

    参数:
    prompt (str): 生成图像的提示。这应该是一个描述性的文本，指导模型生成相应的图像。
    model (str, 可选): 要使用的 DALL-E 模型。默认为 'dall-e-3'。
    n (int, 可选): 要生成的图像数量。默认为 1。
    size (str, 可选): 图像的尺寸，格式为 '宽度x高度'。默认为 '1024x1024'。

    返回:
    dict: 包含生成图像的响应数据。通常包括图像的 URL。
    """
    json_response = GPT.dall_e(prompt, model, n, size)
    result = json_response['data'][0]['url']
    return result


# Google 搜索
def google_search(query):
    """
    根据给定的查询字符串（query），使用 Google 搜索 API 搜索相关信息，并返回搜索结果。

    参数:
    query (str): 用户输入的搜索查询字符串，用于在 Google 搜索中查找相关信息。

    返回:
    dict: 包含 Google 搜索响应的数据。这通常包括搜索结果的标题、链接和简短描述。
    """
    json_response = Google.search(query)

    result_string = ""

    # 检查是否存在 answerBox 并提取 answer
    if 'answerBox' in json_response:
        answer = json_response['answerBox'].get('answer', '')
        if not answer:
            answer = json_response['answerBox'].get('snippet', '')
        result_string += f"标准回答: {answer} \n"

    # 提取并拼接前三个 organic 条目的 snippet
    if 'organic' in json_response:
        organic_list = json_response['organic'][:3] if len(json_response['organic']) > 3 else json_response['organic']
        for i, item in enumerate(organic_list):
            snippet = item['snippet']
            result_string += f"其他回答{i + 1}: {snippet}。\n"
    return result_string
```

等到系统启动后，系统会自动识别 `utils/functions.py` 下的函数，`AutoFunctionGenerator` 会采用 `Few-shot learning`
的方法，高效的为函数生成详细准确的 `JSON Schema`，作为 GPT 所需的 `functions`，并存入数据库 `function-call.functions` 中。

服务对外暴露两个 Restful api：

1. all：用户获取所有可用函数的 `JSON Schema`
2. execute：动态执行函数。直接解析GPT返回的 name 和 arguments，返回字符串（字符串形式不利于扩展，未来会改为JSON）

## 常见问题

1. 如何安装 Docker 和 Docker Compose
    + 查阅官方文档
2. alembic库如何使用？
    + 查阅官方文档

## 捐款

如果觉得这个软件对你有所帮助，欢迎请作者喝咖啡～
