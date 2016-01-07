# libnss_stns
[![Build Status](https://travis-ci.org/pyama86/libnss_stns.svg?branch=master)](https://travis-ci.org/pyama86/libnss_stns)

libnss_stns is [STNS](https://github.com/pyama86/STNS) Client Library.
* /etc/passwd,/etc/groups,/etc/shadow resolver.
* ssh authorized keys wrapper


## install
donload page <https://github.com/pyama86/libnss_stns/releases>
```
$ wget https://github.com/pyama86/libnss_stns/releases/download/<version>/libnss_stns-<version>.noarch.rpm
$ rpm -ivh libnss_stns-<version>.noarch.rpm
```

## config
* /etc/stns/libnss_stns.conf
```
api_end_point = "http://localhost:1104"
```
* /etc/nsswitch.conf
```
~
passwd:     stns files sss ldap
shadow:     stns files sss ldap
group:      stns files sss ldap
~
```

* /etc/sshd/sshd_config
```
~
PubkeyAuthentication yes
AuthorizedKeysCommand /usr/local/bin/ssh_stns_wrapper
AuthorizedKeysCommandUser root
~
```

## tips
advisable to use it together`nscd`(resolver cache service)

## author
* pyama86
