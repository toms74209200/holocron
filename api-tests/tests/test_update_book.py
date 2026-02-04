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


def test_post_books_book_id_with_all_fields_returns_200():
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

    updated_code = random_string()
    updated_title = random_string()
    updated_authors = [random_string(), random_string()]
    updated_publisher = random_string()
    updated_published_date = "2024-03-15"
    updated_thumbnail_url = f"https://example.com/{random_string()}.jpg"

    response = requests.post(
        f"{BASE_URL}/books/{created_book.id}",
        headers={"Authorization": f"Bearer {token}"},
        json={
            "code": updated_code,
            "title": updated_title,
            "authors": updated_authors,
            "publisher": updated_publisher,
            "publishedDate": updated_published_date,
            "thumbnailUrl": updated_thumbnail_url,
        },
    )

    assert response.status_code == 200
    data = response.json()
    assert UUID_PATTERN.match(data["id"])
    assert data["id"] == str(created_book.id)
    assert data["code"] == updated_code
    assert data["title"] == updated_title
    assert data["authors"] == updated_authors
    assert data["publisher"] == updated_publisher
    assert data["publishedDate"] == updated_published_date
    assert data["thumbnailUrl"] == updated_thumbnail_url
    assert data["status"] == "available"
    assert ISO8601_PATTERN.match(data["createdAt"])
    assert ISO8601_PATTERN.match(data["updatedAt"])


def test_post_books_book_id_with_partial_fields_returns_200():
    token = create_user_and_get_token()

    initial_title = random_string()
    initial_authors = [random_string()]
    create_result = post_books.sync_detailed(
        client=AuthenticatedClient(base_url=BASE_URL, token=token),
        body=PostBooksBody(
            title=initial_title,
            authors=initial_authors,
        ),
    )
    assert create_result.status_code == 201
    created_book = create_result.parsed

    updated_title = random_string()

    response = requests.post(
        f"{BASE_URL}/books/{created_book.id}",
        headers={"Authorization": f"Bearer {token}"},
        json={"title": updated_title},
    )

    assert response.status_code == 200
    data = response.json()
    assert data["id"] == str(created_book.id)
    assert data["title"] == updated_title
    assert data["authors"] == initial_authors


def test_post_books_book_id_with_empty_body_returns_200():
    token = create_user_and_get_token()

    initial_title = random_string()
    create_result = post_books.sync_detailed(
        client=AuthenticatedClient(base_url=BASE_URL, token=token),
        body=PostBooksBody(
            title=initial_title,
            authors=[random_string()],
        ),
    )
    assert create_result.status_code == 201
    created_book = create_result.parsed

    response = requests.post(
        f"{BASE_URL}/books/{created_book.id}",
        headers={"Authorization": f"Bearer {token}"},
        json={},
    )

    assert response.status_code == 200
    data = response.json()
    assert data["id"] == str(created_book.id)
    assert data["title"] == initial_title


def test_post_books_book_id_with_invalid_title_returns_400():
    token = create_user_and_get_token()

    create_result = post_books.sync_detailed(
        client=AuthenticatedClient(base_url=BASE_URL, token=token),
        body=PostBooksBody(
            title=random_string(),
            authors=[random_string()],
        ),
    )
    assert create_result.status_code == 201
    created_book = create_result.parsed

    response = requests.post(
        f"{BASE_URL}/books/{created_book.id}",
        headers={"Authorization": f"Bearer {token}"},
        json={"title": ""},
    )

    assert response.status_code == 400
    data = response.json()
    assert data["code"] == "invalid_request"
    assert "message" in data


def test_post_books_book_id_with_too_long_title_returns_400():
    token = create_user_and_get_token()

    create_result = post_books.sync_detailed(
        client=AuthenticatedClient(base_url=BASE_URL, token=token),
        body=PostBooksBody(
            title=random_string(),
            authors=[random_string()],
        ),
    )
    assert create_result.status_code == 201
    created_book = create_result.parsed

    response = requests.post(
        f"{BASE_URL}/books/{created_book.id}",
        headers={"Authorization": f"Bearer {token}"},
        json={"title": "a" * 201},
    )

    assert response.status_code == 400
    data = response.json()
    assert data["code"] == "invalid_request"


def test_post_books_book_id_with_empty_authors_returns_400():
    token = create_user_and_get_token()

    create_result = post_books.sync_detailed(
        client=AuthenticatedClient(base_url=BASE_URL, token=token),
        body=PostBooksBody(
            title=random_string(),
            authors=[random_string()],
        ),
    )
    assert create_result.status_code == 201
    created_book = create_result.parsed

    response = requests.post(
        f"{BASE_URL}/books/{created_book.id}",
        headers={"Authorization": f"Bearer {token}"},
        json={"authors": []},
    )

    assert response.status_code == 400
    data = response.json()
    assert data["code"] == "invalid_request"


def test_post_books_book_id_with_nonexistent_book_returns_404():
    token = create_user_and_get_token()
    nonexistent_id = "00000000-0000-0000-0000-000000000000"

    response = requests.post(
        f"{BASE_URL}/books/{nonexistent_id}",
        headers={"Authorization": f"Bearer {token}"},
        json={"title": random_string()},
    )

    assert response.status_code == 404
    data = response.json()
    assert data["code"] == "not_found"
    assert "message" in data


def test_post_books_book_id_without_auth_returns_401():
    token = create_user_and_get_token()

    create_result = post_books.sync_detailed(
        client=AuthenticatedClient(base_url=BASE_URL, token=token),
        body=PostBooksBody(
            title=random_string(),
            authors=[random_string()],
        ),
    )
    assert create_result.status_code == 201
    created_book = create_result.parsed

    response = requests.post(
        f"{BASE_URL}/books/{created_book.id}",
        json={"title": random_string()},
    )

    assert response.status_code == 401
