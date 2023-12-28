import json
import os

import requests
from dotenv import load_dotenv

load_dotenv('.env')


class Google:
    base_api_url = os.getenv("GOOGLE_API_URL")
    api_key = os.getenv("GOOGLE_API_KEY")
    headers = {
        'Content-Type': 'application/json',
        'X-API-KEY': api_key
    }

    @classmethod
    def _send_request(cls, data):
        try:
            response = requests.post(cls.base_api_url, headers=cls.headers, json=data)
            response.raise_for_status()
            return response.json()
        except requests.RequestException as err:
            print("Request error: ", err)
            return None

    @classmethod
    def search(cls, query):
        data = {"q": query}
        json_response = cls._send_request(data)
        if json_response:
            print(json.dumps(json_response, indent=4, ensure_ascii=False))
        return json_response
