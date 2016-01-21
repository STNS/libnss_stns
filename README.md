# libnss_stns
[![Build Status](https://travis-ci.org/pyama86/libnss_stns.svg?branch=master)](https://travis-ci.org/pyama86/libnss_stns)

libnss_stns is [STNS](https://github.com/pyama86/STNS) Client Library.
* /etc/passwd,/etc/groups,/etc/shadow resolver.
* ssh authorized keys wrapper


## install
download page <https://github.com/pyama86/libnss_stns/releases>
```
$ wget https://github.com/pyama86/libnss_stns/releases/download/<version>/libnss_stns-<version>.noarch.rpm
$ rpm -ivh libnss_stns-<version>.noarch.rpm
```

## config
* /etc/stns/libnss_stns.conf
```
api_end_point = "http://<server>:1104"
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
AuthorizedKeysCommand /usr/local/bin/ssh_stns_wrapper
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

## tips
* auto create home dir by successed ssh login
```
$ echo 'session    required     pam_mkhomedir.so skel=/etc/skel/ umask=0022' >> /etc/pam.d/sshd
```

## author
* pyama86
