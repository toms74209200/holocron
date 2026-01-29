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


def test_post_books_with_valid_input_returns_201(auth_headers):
    response = requests.post(
        f"{BASE_URL}/books",
        json={
            "title": "Test Book",
            "authors": ["Author1", "Author2"],
        },
        headers=auth_headers,
    )

    assert response.status_code == 201
    data = response.json()
    assert UUID_PATTERN.match(data["id"])
    assert data["title"] == "Test Book"
    assert data["authors"] == ["Author1", "Author2"]
    assert data["status"] == "available"
    assert ISO8601_PATTERN.match(data["createdAt"])


def test_post_books_with_optional_fields_returns_201(auth_headers):
    response = requests.post(
        f"{BASE_URL}/books",
        json={
            "title": "Test Book",
            "authors": ["Author1"],
            "publisher": "Test Publisher",
            "publishedDate": "2024-01-01",
            "thumbnailUrl": "https://example.com/thumb.jpg",
        },
        headers=auth_headers,
    )

    assert response.status_code == 201
    data = response.json()
    assert data["title"] == "Test Book"
    assert data["publisher"] == "Test Publisher"
    assert data["publishedDate"] == "2024-01-01"
    assert data["thumbnailUrl"] == "https://example.com/thumb.jpg"


def test_post_books_with_empty_title_returns_400(auth_headers):
    response = requests.post(
        f"{BASE_URL}/books",
        json={
            "title": "",
            "authors": ["Author1"],
        },
        headers=auth_headers,
    )

    assert response.status_code == 400
    data = response.json()
    assert data["code"] == "invalid_request"


def test_post_books_with_empty_authors_returns_400(auth_headers):
    response = requests.post(
        f"{BASE_URL}/books",
        json={
            "title": "Test Book",
            "authors": [],
        },
        headers=auth_headers,
    )

    assert response.status_code == 400
    data = response.json()
    assert data["code"] == "invalid_request"


def test_post_books_with_title_exceeding_200_chars_returns_400(auth_headers):
    response = requests.post(
        f"{BASE_URL}/books",
        json={
            "title": "a" * 201,
            "authors": ["Author1"],
        },
        headers=auth_headers,
    )

    assert response.status_code == 400
    data = response.json()
    assert data["code"] == "invalid_request"


def test_post_books_without_auth_returns_401():
    response = requests.post(
        f"{BASE_URL}/books",
        json={
            "title": "Test Book",
            "authors": ["Author1"],
        },
    )

    assert response.status_code == 401
