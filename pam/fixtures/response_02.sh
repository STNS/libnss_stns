#!/bin/sh
if [ $1 = "sudo/name/example" ]; then
cat << EOS
{
  "metadata": {
    "api_version": 2.0,
    "salt_enable": true,
    "stretching_count": 0,
    "result": "success"
  },
  "items": {
    "example": {
      "id": 1001,
      "password": "",
      "hash_type": "",
      "group_id": 1001,
      "directory": "/home/example",
      "shell": "/bin/sh",
      "gecos": "example",
      "keys": [
        "ssh-rsa xxxxxx"
      ],
      "link_users": null
    }
  }
}
EOS
elif [ $1 = "auth/sudo/name/example/d0255e4d6b4346dafdd47ca17d3f6de91c958ecfffd5f4d08ad32a36f9aa05a9" ]; then
cat << EOS
{
  "metadata": {
    "api_version": 2.0,
    "salt_enable": false,
    "stretching_count": 0,
    "result": "success"
  },
  "items": {
    "example": {
      "id": 1001,
      "password": "",
      "hash_type": "",
      "group_id": 1001,
      "directory": "/home/example",
      "shell": "/bin/sh",
      "gecos": "example",
      "keys": [
        "ssh-rsa xxxxxx"
      ],
      "link_users": null
    }
  }
}
EOS
fi
