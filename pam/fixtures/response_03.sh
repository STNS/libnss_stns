#!/bin/sh
if [ $1 = "sudo/name/example" ]; then
cat << EOS
{
  "metadata": {
    "api_version": 2.0,
    "salt_enable": false,
    "stretching_count": 3,
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
elif [ $1 = "auth/sudo/name/example/5b24f7aa99f1e1da5698a4f91ae0f4b45651a1b625c61ed669dd25ff5b937972" ]; then
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
