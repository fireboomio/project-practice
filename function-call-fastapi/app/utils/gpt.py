import os

import requests
from dotenv import load_dotenv

load_dotenv('.env')


class GPT:
    base_api_url = os.getenv("ONE_API_URL")
    one_api_key = os.getenv("ONE_API_KEY")
    headers = {
        'Content-Type': 'application/json',
        'Authorization': f'Bearer {one_api_key}'
    }

    @classmethod
    def _send_request(cls, endpoint, data):
        url = f"{cls.base_api_url}{endpoint}"
        try:
            response = requests.post(url, headers=cls.headers, json=data)
            response.raise_for_status()
            return response.json()
        except requests.RequestException as err:
            print("Request error: ", err)
            return None

    @classmethod
    def dall_e(cls, prompt, model='dall-e-3', n=1, size='1024x1024'):
        endpoint = "/v1/images/generations"
        data = {
            "model": model or "dall-e-3",
            "prompt": prompt,
            "n": n or 1,
            "size": size or "1024x1024"
        }
        json_response = cls._send_request(endpoint, data)
        if json_response:
            print("JSON Response: ", json_response)
        return json_response

    @classmethod
    def chat(cls, messages, model="gpt-3.5-turbo-1106"):
        endpoint = "/v1/chat/completions"
        data = {
            "model": model,
            "messages": messages,
            "max_tokens": 4096,
        }
        json_response = cls._send_request(endpoint, data)
        if json_response:
            print("JSON Response: ", json_response)
        return json_response
