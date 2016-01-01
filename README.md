# libnss_etcd


## Build
### osx
```
$ brew install docker-machine
$ docker-machine create
$ eval "$(docker-machine env dev)"
$ docker pull golang:1.5
$ docker run -it -v "$(pwd)":/go/src/github.com/pyama86/libnss_etcd \
  -w /go/src/github.com/pyama86/libnss_etcd golang:1.5 make
```
