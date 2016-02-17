#!/bin/bash
eval $(docker-machine env dev)

docker build --no-cache --rm -t stns:libnss_stns . &&
docker run -v "$(pwd)"/binary:/go/src/github.com/STNS/libnss_stns/binary -t stns:libnss_stns

while getopts dr OPT
do
    case $OPT in
        r)
          docker build --no-cache --rm -f docker/rpm -t stns:libnss_stns . && \
          docker run -v "$(pwd)"/binary:/go/src/github.com/STNS/libnss_stns/binary -t stns:libnss_stns
            ;;
        d)
          docker build --no-cache --rm -f docker/deb -t stns:libnss_stns . && \
          docker run -v "$(pwd)"/binary:/go/src/github.com/STNS/libnss_stns/binary -t stns:libnss_stns
            ;;
    esac
done
shift $((OPTIND - 1))
