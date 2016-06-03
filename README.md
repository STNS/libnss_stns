# libnss_stns
[![Build Status](https://travis-ci.org/STNS/libnss_stns.svg?branch=master)](https://travis-ci.org/STNS/libnss_stns)

libnss_stns is [STNS](https://github.com/pyama86/STNS) Client Module.
* /etc/passwd,/etc/groups,/etc/shadow resolver.
* ssh authorized keys wrapper
* sudo auth wrapper

## diagram
![overview](https://cloud.githubusercontent.com/assets/8022082/14065462/37435f26-f465-11e5-9d3a-c25b81cf6ae1.png)

## install
## redhat/centos
```
$ curl -fsSL https://repo.stns.jp/scripts/yum-repo.sh | sh
$ yum install libnss-stns libpam-stns
```
## debian/ubuntu
```
$ curl -fsSL https://repo.stns.jp/scripts/apt-repo.sh | sh
$ apt-get install libnss-stns libpam-stns
```

## upgrade

__Since the cache is funny, nscd and crond want to restart__

## config
* /etc/stns/libnss_stns.conf
```toml
api_end_point = ["http://<server-master>:1104", "http://<server-slave>:1104"]

# support basic auth
user = "basic_user"
password = "basic_password"

# stns access wrapper
wrapper_path = "/usr/local/bin/stns-query-wrapper"

# stns-key-wrapper on failure , this try chain_ssh_wrapper command
chain_ssh_wrapper = "/usr/libexec/openssh/ssh-ldap-wrapper"

ssl_verify = true
```

* /etc/nsswitch.conf
```
passwd:     files stns
shadow:     files stns
group:      files stns
```

* /etc/sshd/sshd_config

```
PubkeyAuthentication yes
AuthorizedKeysCommand /usr/local/bin/stns-key-wrapper
AuthorizedKeysCommandUser root
```

* /etc/nscd.conf
```
        enable-cache            passwd          yes
        positive-time-to-live   passwd          600
        negative-time-to-live   passwd          300
        check-files             passwd          yes
        shared                  group           yes

        enable-cache            group           yes
        positive-time-to-live   group           3600
        negative-time-to-live   group           300
        check-files             group           yes
        shared                  group           yes

        enable-cache            hosts           no
        enable-cache            services        no
        enable-cache            netgroup        no
```

## wrapper command
Query is always through the wrapper.

```
# /usr/local/bin/stns-query-wrapper "user/name/example"
{
  "example": {
    "id": 1001,
    "group_id": 1001,
    "directory": "/home/example",
    "shell": "/bin/sh",
    "gecos": "example",
    "keys": [
      "ssh-rsa xxx
    ],
    "users": null
  }
}
```

Public key is using a key wrapper

```
# /usr/local/bin/stns-key-wrapper "example"
ssh-rsa xxx
```

## Tips
### auto create home dir by successed ssh login
```

$ echo 'session    required     pam_mkhomedir.so skel=/etc/skel/ umask=0022' >> /etc/pam.d/sshd

```

### Sudo Password authentication by Pam

![overview](https://cloud.githubusercontent.com/assets/8022082/14766700/e19c1ef0-0a4e-11e6-975b-1f5cbb4028c8.png)

STNS configure
```toml
[sudoers.example]
password = "sha256hash"
hash_type = "sha256"
```

* /etc/pam.d/sudo
```
#%PAM-1.0
auth       sufficient libpam_stns.so sudo example ⇐
auth       include      system-auth
account    include      system-auth
password   include      system-auth
session    optional     pam_keyinit.so revoke
session    required     pam_limits.soi
...
```
In this case, against STNS to issue query that key example of sudoers, to obtain a password.
This password is the system common.

You will be able to sudo in this only

```
$ sudo ls
[sudo] password for example: <plain password>
example.txt
```

### Password authentication by Pam
STNS attribute
```toml
[users.example]
id = 2000
directory = "/home/exaple"
password = "sha256hash"
hash_type = "sha256"
```

* /etc/pam.d/system-auth or /etc/pam.d/common-auth

```
#%PAM-1.0
# This file is auto-generated.
# User changes will be destroyed the next time authconfig is run.
auth        required      pam_env.so
auth        sufficient    pam_unix.so nullok try_first_pass
auth        sufficient    libpam_stns.so ⇐
```

In this way you can make the password authentication in the system-wide by configure

# develop test
Use docker machine

```
$ docker-machine start dev
$ eval $(docker-machine env dev)
$ rake test
```

## author
* pyama86
