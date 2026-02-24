import re
from datetime import datetime

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


def test_get_users_me_borrowings_with_no_borrowed_books_returns_200():
    token = create_user_and_get_token()

    response = requests.get(
        f"{BASE_URL}/users/me/borrowings",
        headers={"Authorization": f"Bearer {token}"},
    )

    assert response.status_code == 200
    data = response.json()
    assert data["items"] == []
    assert data["total"] == 0


def test_get_users_me_borrowing_with_borrowed_book_returns_book():
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
    borrow_data = borrow_response.json()

    response = requests.get(
        f"{BASE_URL}/users/me/borrowings",
        headers={"Authorization": f"Bearer {token}"},
    )

    assert response.status_code == 200
    data = response.json()
    assert data["total"] == 1
    assert len(data["items"]) == 1

    item = data["items"][0]
    assert UUID_PATTERN.match(item["id"])
    assert item["id"] == str(created_book.id)
    assert item["title"] == unique_title
    assert isinstance(item["authors"], list)
    assert len(item["authors"]) > 0
    assert ISO8601_PATTERN.match(item["borrowedAt"])
    assert ISO8601_PATTERN.match(item["dueDate"])
    assert datetime.fromisoformat(
        item["borrowedAt"].replace("Z", "+00:00")
    ) == datetime.fromisoformat(borrow_data["borrowedAt"].replace("Z", "+00:00"))
    assert datetime.fromisoformat(
        item["dueDate"].replace("Z", "+00:00")
    ) == datetime.fromisoformat(borrow_data["dueDate"].replace("Z", "+00:00"))


def test_get_users_me_borrowing_after_return_shows_empty():
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

    return_response = requests.post(
        f"{BASE_URL}/books/{created_book.id}/return",
        headers={"Authorization": f"Bearer {token}"},
    )
    assert return_response.status_code == 200

    response = requests.get(
        f"{BASE_URL}/users/me/borrowings",
        headers={"Authorization": f"Bearer {token}"},
    )

    assert response.status_code == 200
    data = response.json()
    assert data["items"] == []
    assert data["total"] == 0


def test_get_users_me_borrowing_only_shows_own_books():
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

    response = requests.get(
        f"{BASE_URL}/users/me/borrowings",
        headers={"Authorization": f"Bearer {token2}"},
    )

    assert response.status_code == 200
    data = response.json()
    assert data["items"] == []
    assert data["total"] == 0


def test_get_users_me_borrowing_without_auth_returns_401():
    response = requests.get(
        f"{BASE_URL}/users/me/borrowings",
    )

    assert response.status_code == 401
