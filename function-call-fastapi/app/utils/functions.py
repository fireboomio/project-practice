from app.utils.google import Google
from app.utils.gpt import GPT


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
    response = GPT.dall_e(prompt, model, n, size)
    return response


def google_search(query):
    """
    根据给定的查询字符串（query），使用 Google 搜索 API 搜索相关信息，并返回搜索结果。

    参数:
    query (str): 用户输入的搜索查询字符串，用于在 Google 搜索中查找相关信息。

    返回:
    dict: 包含 Google 搜索响应的数据。这通常包括搜索结果的标题、链接和简短描述。
    """
    response = Google.search(query)
    return response
