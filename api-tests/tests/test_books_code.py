import re

import pytest
import requests

from lib.api_config import BASE_URL
from lib.auth import create_user_and_get_token

UUID_PATTERN = re.compile(
    r"^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"
)
ISO8601_PATTERN = re.compile(
    r"^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(Z|[+-]\d{2}:\d{2})$"
)


@pytest.fixture(scope="module")
def auth_headers():
    token = create_user_and_get_token()
    return {"Authorization": f"Bearer {token}"}


def test_post_books_code_with_valid_isbn_returns_201(auth_headers):
    response = requests.post(
        f"{BASE_URL}/books/code",
        json={"code": "9784873115658"},
        headers=auth_headers,
    )

    assert response.status_code == 201
    data = response.json()
    assert UUID_PATTERN.match(data["id"])
    assert data["code"] == "9784873115658"
    assert data["title"] == "リーダブルコード"
    assert isinstance(data["authors"], list)
    assert len(data["authors"]) > 0
    assert data["status"] == "available"
    assert ISO8601_PATTERN.match(data["createdAt"])


def test_post_books_code_with_empty_code_returns_400(auth_headers):
    response = requests.post(
        f"{BASE_URL}/books/code",
        json={"code": ""},
        headers=auth_headers,
    )

    assert response.status_code == 400
    data = response.json()
    assert data["code"] == "invalid_request"


def test_post_books_code_without_auth_returns_401():
    response = requests.post(
        f"{BASE_URL}/books/code",
        json={"code": "9784873115658"},
    )

    assert response.status_code == 401
