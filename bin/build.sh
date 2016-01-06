eval $(docker-machine env dev)
docker build -t centos:libnss . && docker run -v "$(pwd)"/releases:/go/src/github.com/pyama86/libnss_stns/RPM/RPMS centos:libnss
