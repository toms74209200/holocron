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


def test_get_books_returns_registered_book():
    token = create_user_and_get_token()

    unique_title = random_string()
    result = post_books.sync_detailed(
        client=AuthenticatedClient(base_url=BASE_URL, token=token),
        body=PostBooksBody(title=unique_title, authors=[random_string()]),
    )
    assert result.status_code == 201
    created_book = result.parsed

    response = requests.get(
        f"{BASE_URL}/books",
        headers={"Authorization": f"Bearer {token}"},
    )

    assert response.status_code == 200
    data = response.json()
    assert data["total"] >= 1

    book_ids = [item["id"] for item in data["items"]]
    assert str(created_book.id) in book_ids
    if len(data["items"]) > 0:
        item = data["items"][0]

        assert "id" in item
        assert "title" in item
        assert "authors" in item
        assert "status" in item
        assert "createdAt" in item

        assert UUID_PATTERN.match(item["id"])
        assert isinstance(item["authors"], list)
        assert item["status"] in ["available", "borrowed"]
        assert ISO8601_PATTERN.match(item["createdAt"])


def test_get_books_with_limit_parameter():
    token = create_user_and_get_token()

    unique_title = random_string()
    result = post_books.sync_detailed(
        client=AuthenticatedClient(base_url=BASE_URL, token=token),
        body=PostBooksBody(title=unique_title, authors=[random_string()]),
    )
    assert result.status_code == 201

    response = requests.get(
        f"{BASE_URL}/books",
        params={"limit": 5},
        headers={"Authorization": f"Bearer {token}"},
    )

    assert response.status_code == 200
    data = response.json()
    assert len(data["items"]) <= 5


def test_get_books_with_offset_parameter():
    token = create_user_and_get_token()

    unique_title = random_string()
    result = post_books.sync_detailed(
        client=AuthenticatedClient(base_url=BASE_URL, token=token),
        body=PostBooksBody(title=unique_title, authors=[random_string()]),
    )
    assert result.status_code == 201

    response = requests.get(
        f"{BASE_URL}/books",
        headers={"Authorization": f"Bearer {token}"},
    )
    assert response.status_code == 200
    total = response.json()["total"]

    response_with_offset = requests.get(
        f"{BASE_URL}/books",
        params={"offset": 1},
        headers={"Authorization": f"Bearer {token}"},
    )

    assert response_with_offset.status_code == 200
    data = response_with_offset.json()

    assert data["total"] == total


def test_get_books_with_search_query():
    token = create_user_and_get_token()

    unique_title = random_string()
    result = post_books.sync_detailed(
        client=AuthenticatedClient(base_url=BASE_URL, token=token),
        body=PostBooksBody(title=unique_title, authors=[random_string()]),
    )
    assert result.status_code == 201

    response = requests.get(
        f"{BASE_URL}/books",
        params={"q": unique_title},
        headers={"Authorization": f"Bearer {token}"},
    )

    assert response.status_code == 200
    data = response.json()
    assert data["total"] >= 1
    titles = [item["title"] for item in data["items"]]
    assert unique_title in titles


def test_get_books_without_auth_returns_401():
    response = requests.get(f"{BASE_URL}/books")

    assert response.status_code == 401


def test_get_books_with_invalid_limit_returns_400():
    token = create_user_and_get_token()
    response = requests.get(
        f"{BASE_URL}/books",
        params={"limit": "invalid"},
        headers={"Authorization": f"Bearer {token}"},
    )

    assert response.status_code == 400


def test_get_books_with_invalid_offset_returns_400():
    token = create_user_and_get_token()
    response = requests.get(
        f"{BASE_URL}/books",
        params={"offset": "invalid"},
        headers={"Authorization": f"Bearer {token}"},
    )

    assert response.status_code == 400


def test_get_books_after_borrow_returns_borrower_info():
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
        f"{BASE_URL}/books",
        headers={"Authorization": f"Bearer {token}"},
    )

    assert response.status_code == 200
    data = response.json()

    borrowed_book = next(
        (item for item in data["items"] if item["id"] == str(created_book.id)),
        None
    )
    assert borrowed_book is not None
    assert borrowed_book["status"] == "borrowed"
    assert "borrower" in borrowed_book
    assert borrowed_book["borrower"] is not None
    assert UUID_PATTERN.match(borrowed_book["borrower"]["id"])
    assert isinstance(borrowed_book["borrower"]["name"], str)
    assert ISO8601_PATTERN.match(borrowed_book["borrower"]["borrowedAt"])
