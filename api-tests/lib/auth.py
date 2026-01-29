import os

import requests

from lib.api_config import BASE_URL

FIREBASE_API_KEY = os.getenv("FIREBASE_API_KEY", "fake-api-key")
FIREBASE_AUTH_EMULATOR_HOST = os.getenv("FIREBASE_AUTH_EMULATOR_HOST", "firebase:9099")
FIREBASE_AUTH_URL = f"http://{FIREBASE_AUTH_EMULATOR_HOST}/identitytoolkit.googleapis.com/v1/accounts:signInWithCustomToken"


def create_user_and_get_token(name: str | None = None) -> str:
    payload = {"name": name} if name else {}
    response = requests.post(f"{BASE_URL}/users", json=payload)
    response.raise_for_status()
    custom_token = response.json()["customToken"]

    if not FIREBASE_API_KEY:
        raise ValueError("FIREBASE_API_KEY environment variable is not set")

    response = requests.post(
        f"{FIREBASE_AUTH_URL}?key={FIREBASE_API_KEY}",
        json={
            "token": custom_token,
            "returnSecureToken": True,
        },
    )
    response.raise_for_status()
    return response.json()["idToken"]
