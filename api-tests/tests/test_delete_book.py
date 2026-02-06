import requests

from lib.api_config import BASE_URL
from lib.auth import create_user_and_get_token
from lib.random_string import random_string
from openapi_gen.holocron_library_management_api_client import AuthenticatedClient
from openapi_gen.holocron_library_management_api_client.api.books import post_books
from openapi_gen.holocron_library_management_api_client.models import PostBooksBody


def test_delete_book_with_available_book_returns_204():
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

    response = requests.delete(
        f"{BASE_URL}/books/{created_book.id}",
        headers={"Authorization": f"Bearer {token}"},
        json={"reason": "disposal"},
    )

    assert response.status_code == 204
    assert response.text == ""

    get_response = requests.get(
        f"{BASE_URL}/books/{created_book.id}",
        headers={"Authorization": f"Bearer {token}"},
    )
    assert get_response.status_code == 404


def test_delete_book_with_nonexistent_book_returns_404():
    token = create_user_and_get_token()
    nonexistent_id = "00000000-0000-0000-0000-000000000000"

    response = requests.delete(
        f"{BASE_URL}/books/{nonexistent_id}",
        headers={"Authorization": f"Bearer {token}"},
        json={"reason": "disposal"},
    )

    assert response.status_code == 404
    data = response.json()
    assert data["code"] == "not_found"
    assert "message" in data


def test_delete_book_with_borrowed_book_returns_409():
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

    borrow_response = requests.post(
        f"{BASE_URL}/books/{created_book.id}/borrow",
        headers={"Authorization": f"Bearer {token}"},
    )
    assert borrow_response.status_code == 200

    response = requests.delete(
        f"{BASE_URL}/books/{created_book.id}",
        headers={"Authorization": f"Bearer {token}"},
        json={"reason": "disposal"},
    )

    assert response.status_code == 409
    data = response.json()
    assert data["code"] == "conflict"
    assert "message" in data


def test_delete_book_with_returned_book_returns_204():
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

    response = requests.delete(
        f"{BASE_URL}/books/{created_book.id}",
        headers={"Authorization": f"Bearer {token}"},
        json={"reason": "transfer", "memo": "友人に譲りました"},
    )

    assert response.status_code == 204


def test_delete_book_with_invalid_reason_returns_400():
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

    response = requests.delete(
        f"{BASE_URL}/books/{created_book.id}",
        headers={"Authorization": f"Bearer {token}"},
        json={"reason": "invalid_reason"},
    )

    assert response.status_code == 400
    data = response.json()
    assert data["code"] == "invalid_request"


def test_delete_book_without_reason_returns_400():
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

    response = requests.delete(
        f"{BASE_URL}/books/{created_book.id}",
        headers={"Authorization": f"Bearer {token}"},
        json={},
    )

    assert response.status_code == 400


def test_delete_book_without_auth_returns_401():
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

    response = requests.delete(
        f"{BASE_URL}/books/{created_book.id}",
        json={"reason": "disposal"},
    )

    assert response.status_code == 401
