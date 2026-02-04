import re

import requests

from lib.api_config import BASE_URL
from lib.auth import create_user_and_get_token
from lib.random_string import random_string
from openapi_gen.holocron_library_management_api_client import AuthenticatedClient
from openapi_gen.holocron_library_management_api_client.api.books import post_books
from openapi_gen.holocron_library_management_api_client.models import PostBooksBody

UUID_PATTERN = re.compile(
    r"^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"
)
ISO8601_PATTERN = re.compile(
    r"^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(Z|[+-]\d{2}:\d{2})$"
)


def test_get_book_with_valid_id_returns_200():
    token = create_user_and_get_token()

    unique_title = random_string()
    author_name = random_string()
    create_result = post_books.sync_detailed(
        client=AuthenticatedClient(base_url=BASE_URL, token=token),
        body=PostBooksBody(
            title=unique_title,
            authors=[author_name],
            publisher="Test Publisher",
            published_date="2024-01-01",
            thumbnail_url="https://example.com/thumb.jpg",
        ),
    )
    assert create_result.status_code == 201
    created_book = create_result.parsed

    response = requests.get(
        f"{BASE_URL}/books/{created_book.id}",
        headers={"Authorization": f"Bearer {token}"},
    )

    assert response.status_code == 200
    data = response.json()
    assert UUID_PATTERN.match(data["id"])
    assert data["id"] == str(created_book.id)
    assert data["title"] == unique_title
    assert data["authors"] == [author_name]
    assert data["publisher"] == "Test Publisher"
    assert data["publishedDate"] == "2024-01-01"
    assert data["thumbnailUrl"] == "https://example.com/thumb.jpg"
    assert data["status"] == "available"
    assert ISO8601_PATTERN.match(data["createdAt"])


def test_get_book_with_nonexistent_id_returns_404():
    token = create_user_and_get_token()
    nonexistent_id = "00000000-0000-0000-0000-000000000000"

    response = requests.get(
        f"{BASE_URL}/books/{nonexistent_id}",
        headers={"Authorization": f"Bearer {token}"},
    )

    assert response.status_code == 404
    data = response.json()
    assert data["code"] == "not_found"
    assert "message" in data


def test_get_book_without_auth_returns_401():
    token = create_user_and_get_token()

    unique_title = random_string()
    create_result = post_books.sync_detailed(
        client=AuthenticatedClient(base_url=BASE_URL, token=token),
        body=PostBooksBody(title=unique_title, authors=[random_string()]),
    )
    assert create_result.status_code == 201
    created_book = create_result.parsed

    response = requests.get(
        f"{BASE_URL}/books/{created_book.id}",
    )

    assert response.status_code == 401


def test_get_book_after_borrow_returns_borrowed_status():
    token = create_user_and_get_token()

    unique_title = random_string()
    create_result = post_books.sync_detailed(
        client=AuthenticatedClient(base_url=BASE_URL, token=token),
        body=PostBooksBody(title=unique_title, authors=[random_string()]),
    )
    assert create_result.status_code == 201
    created_book = create_result.parsed

    borrow_response = requests.post(
        f"{BASE_URL}/books/{created_book.id}/borrow",
        headers={"Authorization": f"Bearer {token}"},
    )
    assert borrow_response.status_code == 200

    response = requests.get(
        f"{BASE_URL}/books/{created_book.id}",
        headers={"Authorization": f"Bearer {token}"},
    )

    assert response.status_code == 200
    data = response.json()
    assert data["status"] == "borrowed"
    assert "borrower" in data
    assert data["borrower"] is not None
    assert UUID_PATTERN.match(data["borrower"]["id"])
    assert isinstance(data["borrower"]["name"], str)
    assert ISO8601_PATTERN.match(data["borrower"]["borrowedAt"])
