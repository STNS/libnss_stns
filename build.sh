#!/bin/bash
eval $(docker-machine env dev)

while getopts brd OPT
do
    case $OPT in
        b)
          docker build --no-cache --rm -t stns:libnss_stns .
            ;;
        r)
          docker build --no-cache --rm -f docker/rpm -t stns:libnss_stns .
            ;;
        d)
          docker build --no-cache --rm -f docker/deb -t stns:libnss_stns .
            ;;
    esac
done
shift $((OPTIND - 1))
docker run -v "$(pwd)"/binary:/go/src/github.com/STNS/libnss_stns/binary -t stns:libnss_stns

