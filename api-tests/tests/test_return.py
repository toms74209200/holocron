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


def test_post_books_return_with_borrowed_book_returns_200():
    token = create_user_and_get_token()

    unique_title = random_string()
    create_result = post_books.sync_detailed(
        client=AuthenticatedClient(base_url=BASE_URL, token=token),
        body=PostBooksBody(
            title=unique_title,
            authors=[random_string()],
        ),
    )
    assert create_result.status_code == 201
    created_book = create_result.parsed

    borrow_response = requests.post(
        f"{BASE_URL}/books/{created_book.id}/borrow",
        headers={"Authorization": f"Bearer {token}"},
    )
    assert borrow_response.status_code == 200

    response = requests.post(
        f"{BASE_URL}/books/{created_book.id}/return",
        headers={"Authorization": f"Bearer {token}"},
    )

    assert response.status_code == 200
    data = response.json()
    assert UUID_PATTERN.match(data["id"])
    assert data["id"] == str(created_book.id)
    assert data["title"] == unique_title
    assert isinstance(data["authors"], list)
    assert len(data["authors"]) > 0
    assert data["status"] == "available"
    assert "borrower" not in data
    assert ISO8601_PATTERN.match(data["createdAt"])


def test_post_books_return_with_not_borrowed_book_returns_409():
    token = create_user_and_get_token()

    unique_title = random_string()
    create_result = post_books.sync_detailed(
        client=AuthenticatedClient(base_url=BASE_URL, token=token),
        body=PostBooksBody(
            title=unique_title,
            authors=[random_string()],
        ),
    )
    assert create_result.status_code == 201
    created_book = create_result.parsed

    response = requests.post(
        f"{BASE_URL}/books/{created_book.id}/return",
        headers={"Authorization": f"Bearer {token}"},
    )

    assert response.status_code == 409
    data = response.json()
    assert data["code"] == "not_borrowed"
    assert "message" in data


def test_post_books_return_with_nonexistent_book_returns_404():
    token = create_user_and_get_token()
    nonexistent_id = "00000000-0000-0000-0000-000000000000"

    response = requests.post(
        f"{BASE_URL}/books/{nonexistent_id}/return",
        headers={"Authorization": f"Bearer {token}"},
    )

    assert response.status_code == 404
    data = response.json()
    assert data["code"] == "book_not_found"
    assert "message" in data


def test_post_books_return_without_auth_returns_401():
    token = create_user_and_get_token()

    unique_title = random_string()
    create_result = post_books.sync_detailed(
        client=AuthenticatedClient(base_url=BASE_URL, token=token),
        body=PostBooksBody(
            title=unique_title,
            authors=[random_string()],
        ),
    )
    assert create_result.status_code == 201
    created_book = create_result.parsed

    borrow_response = requests.post(
        f"{BASE_URL}/books/{created_book.id}/borrow",
        headers={"Authorization": f"Bearer {token}"},
    )
    assert borrow_response.status_code == 200

    response = requests.post(
        f"{BASE_URL}/books/{created_book.id}/return",
    )

    assert response.status_code == 401


def test_post_books_return_by_different_user_returns_403():
    token1 = create_user_and_get_token()
    token2 = create_user_and_get_token()

    unique_title = random_string()
    create_result = post_books.sync_detailed(
        client=AuthenticatedClient(base_url=BASE_URL, token=token1),
        body=PostBooksBody(
            title=unique_title,
            authors=[random_string()],
        ),
    )
    assert create_result.status_code == 201
    created_book = create_result.parsed

    borrow_response = requests.post(
        f"{BASE_URL}/books/{created_book.id}/borrow",
        headers={"Authorization": f"Bearer {token1}"},
    )
    assert borrow_response.status_code == 200

    response = requests.post(
        f"{BASE_URL}/books/{created_book.id}/return",
        headers={"Authorization": f"Bearer {token2}"},
    )

    assert response.status_code == 403
    data = response.json()
    assert data["code"] == "forbidden"
    assert "message" in data


def test_get_book_after_return_shows_available_status():
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

    return_response = requests.post(
        f"{BASE_URL}/books/{created_book.id}/return",
        headers={"Authorization": f"Bearer {token}"},
    )
    assert return_response.status_code == 200

    response = requests.get(
        f"{BASE_URL}/books/{created_book.id}",
        headers={"Authorization": f"Bearer {token}"},
    )

    assert response.status_code == 200
    data = response.json()
    assert data["status"] == "available"
    assert data.get("borrower") is None
