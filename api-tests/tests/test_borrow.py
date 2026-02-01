import re
from datetime import datetime, timedelta

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


def test_post_books_borrow_with_available_book_returns_200():
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
        f"{BASE_URL}/books/{created_book.id}/borrow",
        headers={"Authorization": f"Bearer {token}"},
    )

    assert response.status_code == 200
    data = response.json()
    assert UUID_PATTERN.match(data["id"])
    assert UUID_PATTERN.match(data["bookId"])
    assert data["bookId"] == str(created_book.id)
    assert UUID_PATTERN.match(data["borrowerId"])
    assert ISO8601_PATTERN.match(data["borrowedAt"])
    assert ISO8601_PATTERN.match(data["dueDate"])

    borrowed_at = datetime.fromisoformat(data["borrowedAt"].replace("Z", "+00:00"))
    due_date = datetime.fromisoformat(data["dueDate"].replace("Z", "+00:00"))
    expected_due_date = borrowed_at + timedelta(days=7)
    assert due_date.date() == expected_due_date.date()


def test_post_books_borrow_with_custom_due_days_returns_200():
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
        f"{BASE_URL}/books/{created_book.id}/borrow",
        headers={"Authorization": f"Bearer {token}"},
        json={"dueDays": 14},
    )

    assert response.status_code == 200
    data = response.json()
    assert UUID_PATTERN.match(data["id"])
    assert data["bookId"] == str(created_book.id)
    assert ISO8601_PATTERN.match(data["dueDate"])

    borrowed_at = datetime.fromisoformat(data["borrowedAt"].replace("Z", "+00:00"))
    due_date = datetime.fromisoformat(data["dueDate"].replace("Z", "+00:00"))
    expected_due_date = borrowed_at + timedelta(days=14)
    assert due_date.date() == expected_due_date.date()


def test_post_books_borrow_same_user_extends_due_date():
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

    first_borrow_response = requests.post(
        f"{BASE_URL}/books/{created_book.id}/borrow",
        headers={"Authorization": f"Bearer {token}"},
    )
    assert first_borrow_response.status_code == 200
    first_data = first_borrow_response.json()
    first_due_date = datetime.fromisoformat(first_data["dueDate"].replace("Z", "+00:00"))

    response = requests.post(
        f"{BASE_URL}/books/{created_book.id}/borrow",
        headers={"Authorization": f"Bearer {token}"},
    )

    assert response.status_code == 200
    data = response.json()
    assert UUID_PATTERN.match(data["id"])
    assert data["bookId"] == str(created_book.id)
    assert ISO8601_PATTERN.match(data["dueDate"])

    new_due_date = datetime.fromisoformat(data["dueDate"].replace("Z", "+00:00"))
    expected_extended_due_date = first_due_date + timedelta(days=7)
    assert new_due_date.date() == expected_extended_due_date.date()


def test_post_books_borrow_with_already_borrowed_book_returns_409():
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

    first_borrow_response = requests.post(
        f"{BASE_URL}/books/{created_book.id}/borrow",
        headers={"Authorization": f"Bearer {token1}"},
    )
    assert first_borrow_response.status_code == 200

    response = requests.post(
        f"{BASE_URL}/books/{created_book.id}/borrow",
        headers={"Authorization": f"Bearer {token2}"},
    )

    assert response.status_code == 409
    data = response.json()
    assert data["code"] == "book_already_borrowed"
    assert "message" in data


def test_post_books_borrow_with_invalid_due_days_returns_400():
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
        f"{BASE_URL}/books/{created_book.id}/borrow",
        headers={"Authorization": f"Bearer {token}"},
        json={"dueDays": 0},
    )

    assert response.status_code == 400
    data = response.json()
    assert data["code"] == "invalid_request"
    assert "message" in data


def test_post_books_borrow_with_nonexistent_book_returns_404():
    token = create_user_and_get_token()
    nonexistent_id = "00000000-0000-0000-0000-000000000000"

    response = requests.post(
        f"{BASE_URL}/books/{nonexistent_id}/borrow",
        headers={"Authorization": f"Bearer {token}"},
    )

    assert response.status_code == 404
    data = response.json()
    assert data["code"] == "book_not_found"
    assert "message" in data


def test_post_books_borrow_without_auth_returns_401():
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
        f"{BASE_URL}/books/{created_book.id}/borrow",
    )

    assert response.status_code == 401
