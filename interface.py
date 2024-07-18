#!/usr/bin/env python3

import click
import requests
import filetype
from requests_toolbelt.multipart.encoder import MultipartEncoder

BASE_URL = "http://localhost:3000"  # Replace with your actual base URL
AUTH = f"{BASE_URL}/auth"
API_URL = f"{BASE_URL}/api/v1"
TOKEN_FILE = ".token"


def set_token(token: str) -> None:
    with open(TOKEN_FILE, "w") as f:
        f.write(token)


def get_token() -> str:
    try:
        with open(TOKEN_FILE, "r") as f:
            return f.read().strip()
    except FileNotFoundError:
        return None


@click.group()
def cli() -> None:
    pass


@cli.command()
@click.option("--name", prompt=True, help="The name of the user")
@click.option(
    "--password", prompt=True, hide_input=True, help="The password of the user"
)
def register(name: str, password: str) -> None:
    url = f"{AUTH}/register"
    fields = {
        "name": name,
        "password": password,
    }
    m = MultipartEncoder(fields=fields)
    headers = {"Content-Type": m.content_type}

    response = requests.post(url, data=m, headers=headers)
    if response.status_code == 201:
        click.echo("User registered successfully")
    else:
        click.echo(f"Failed to register user: {response.json().get('error')}")


@cli.command()
@click.option("--name", prompt=True, help="The name of the user")
@click.option(
    "--password", prompt=True, hide_input=True, help="The password of the user"
)
def login(
    name: str,
    password: str,
) -> None:
    url = f"{AUTH}/login"
    payload = {"name": name, "password": password}
    response = requests.post(url, json=payload)
    if response.status_code == 200:
        token = response.json().get("token")
        set_token(token)
        click.echo(f"Logged in successfully, token: {token}")
    else:
        click.echo(f"Failed to login: {response.json().get('error')}")


@cli.command()
@click.option("--token", help="The JWT token of the user")
def get_users(token: str) -> None:
    token = token or get_token()
    if not token:
        click.echo("No token provided or found. Please login first.")
        return
    url = f"{API_URL}/users"
    headers = {"Authorization": f"Bearer {token}"}
    response = requests.get(url, headers=headers)
    if response.status_code == 200:
        users = response.json().get("users")
        for user in users:
            click.echo(f"User ID: {user['id']}, Name: {user['name']}")
    else:
        click.echo(f"Failed to fetch users: {response.json().get('error')}")


@cli.command()
@click.option("--token", help="The JWT token of the user")
@click.option("--content", prompt=True, help="The content of the post")
def create_post(
    token: str,
    content: str,
) -> None:
    token = token or get_token()
    if not token:
        click.echo("No token provided or found. Please login first.")
        return
    url = f"{API_URL}/posts"
    headers = {"Authorization": f"Bearer {token}"}
    payload = {"content": content}
    response = requests.post(url, json=payload, headers=headers)
    if response.status_code == 201:
        click.echo("Post created successfully")
    else:
        click.echo(f"Failed to create post: {response.json().get('error')}")


@cli.command()
@click.option("--token", help="The JWT token of the user")
@click.option("--id", prompt=True, help="The ID of the post to like")
def like_post(
    token: str,
    id: str,
) -> None:
    token = token or get_token()
    if not token:
        click.echo("No token provided or found. Please login first.")
        return
    url = f"{API_URL}/posts/{id}/like"
    headers = {"Authorization": f"Bearer {token}"}
    response = requests.post(url, headers=headers)
    if response.status_code == 200:
        click.echo("Post liked successfully")
    else:
        click.echo(f"Failed to like post: {response.json().get('error')}")


@cli.command()
@click.option("--token", help="The JWT token of the user")
def read_notifications(token: str) -> None:
    token = token or get_token()
    if not token:
        click.echo("No token provided or found. Please login first.")
        return
    url = f"{API_URL}/notifications"
    headers = {"Authorization": f"Bearer {token}"}
    response = requests.get(url, headers=headers)
    if response.status_code == 200:
        notifications = response.json().get("notifications")
        for notification in notifications:
            click.echo(
                f"Notification ID: {notification['id']}, Type: {notification['type']}, Post ID: {notification['post_id']}, Liked By: {notification['liked_by']}"
            )
    else:
        click.echo(f"Failed to fetch notifications: {response.json().get('error')}")


@cli.command()
@click.option("--token", help="The JWT token of the user")
@click.option("--user_id", prompt=True, help="The ID of the user to update")
@click.option("--name", help="The new name of the user")
@click.option(
    "--avatar", type=click.Path(exists=True), help="The path to the new avatar image"
)
def update_user(
    token: str,
    user_id: str,
    name: str,
    avatar: str,
) -> None:
    token = token or get_token()
    if not token:
        click.echo("No token provided or found. Please login first.")
        return
    url = f"{API_URL}/users/{user_id}"
    headers = {"Authorization": f"Bearer {token}"}

    fields = {}
    if name:
        fields["name"] = name
    if avatar:
        with open(avatar, "rb") as avatar_file:
            kind = filetype.guess(avatar_file)
            if kind is None:
                click.echo("Cannot determine file type of the avatar")
                return
            avatar_type = kind.mime
            avatar_file.seek(0)  # Reset file pointer to the beginning
            fields["avatar"] = (avatar, avatar_file, avatar_type)

    m = MultipartEncoder(fields=fields)
    headers["Content-Type"] = m.content_type

    response = requests.patch(url, data=m, headers=headers)
    if response.status_code == 200:
        click.echo("User updated successfully")
    else:
        click.echo(f"Failed to update user: {response.json().get('error')}")


if __name__ == "__main__":
    cli()
