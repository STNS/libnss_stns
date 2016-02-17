# libnss_stns
[![Build Status](https://travis-ci.org/STNS/libnss_stns.svg?branch=master)](https://travis-ci.org/STNS/libnss_stns)

libnss_stns is [STNS](https://github.com/pyama86/STNS) Client Library.
* /etc/passwd,/etc/groups,/etc/shadow resolver.
* ssh authorized keys wrapper

## install
## redhat/centos
```
$ curl -fsSL https://repo.stns.jp/scripts/yum-repo.sh | sh
$ yum install libnss-stns
```
## debian/ubuntu
```
$ curl -fsSL https://repo.stns.jp/scripts/apt-repo.sh | sh
$ apt-get install libnss-stns
```

## config
* /etc/stns/libnss_stns.conf
```
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
~
passwd:     files stns
shadow:     files stns
group:      files stns
~
```

* /etc/sshd/sshd_config

```
~
PubkeyAuthentication yes
AuthorizedKeysCommand /usr/local/bin/stns-key-wrapper
AuthorizedKeysCommandUser root
~
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
* stns-query-wrapper "user/name/pyama"
stns http query wrapper
```
# /usr/local/bin/stns-query-wrapper "user/name/pyama"
{
  "pyama": {
    "id": 1001,
    "group_id": 1001,
    "directory": "/home/pyama",
    "shell": "/bin/sh",
    "gecos": "pyama",
    "keys": [
      "ssh-rsa xxx
    ],
    "users": null
  }
}
```

* stns-key-wrapper "pyama"
stns ssh key wrapper
```
# /usr/local/bin/stns-key-wrapper "pyama"
ssh-rsa xxx
```
## tips
* auto create home dir by successed ssh login
```
$ echo 'session    required     pam_mkhomedir.so skel=/etc/skel/ umask=0022' >> /etc/pam.d/sshd
```

## author
* pyama86
