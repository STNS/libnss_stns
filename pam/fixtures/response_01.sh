#!/bin/sh
if [ $1 = "sudo/name/example" ]; then
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
elif [ $1 = "auth/sudo/name/example/9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08" ]; then
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
