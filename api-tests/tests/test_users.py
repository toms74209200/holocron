import re

import requests

from lib.api_config import BASE_URL

UUID_PATTERN = re.compile(
    r"^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"
)
ISO8601_PATTERN = re.compile(
    r"^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(Z|[+-]\d{2}:\d{2})$"
)


def test_post_users_with_valid_name_returns_201():
    response = requests.post(
        f"{BASE_URL}/users",
        json={"name": "TestUser"},
    )

    assert response.status_code == 201
    data = response.json()
    assert UUID_PATTERN.match(data["id"])
    assert data["name"] == "TestUser"
    assert "customToken" in data
    assert ISO8601_PATTERN.match(data["createdAt"])


def test_post_users_with_empty_name_returns_201_with_generated_name():
    response = requests.post(
        f"{BASE_URL}/users",
        json={"name": ""},
    )

    assert response.status_code == 201
    data = response.json()
    assert UUID_PATTERN.match(data["id"])
    assert data["name"].startswith("ユーザー")
    assert "customToken" in data


def test_post_users_with_no_name_returns_201_with_generated_name():
    response = requests.post(
        f"{BASE_URL}/users",
        json={},
    )

    assert response.status_code == 201
    data = response.json()
    assert UUID_PATTERN.match(data["id"])
    assert data["name"].startswith("ユーザー")
    assert "customToken" in data


def test_post_users_with_name_exceeding_50_chars_returns_400():
    response = requests.post(
        f"{BASE_URL}/users",
        json={"name": "a" * 51},
    )

    assert response.status_code == 400
    data = response.json()
    assert data["code"] == "invalid_request"
