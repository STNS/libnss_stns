#!/bin/bash
eval $(docker-machine env dev)
rm -rf ./binary/*

docker build --no-cache --rm -t stns:libnss_stns . &&
docker run -v "$(pwd)"/binary:/go/src/github.com/STNS/libnss_stns/binary -t stns:libnss_stns

func_ssh_build()
{
  cd ../ssh_stns_wrapper
  ./build.sh
  cp binary/* ../libnss_stns/binary
  cd ../libnss_stns
}

while getopts dr OPT
do
    case $OPT in
        r)
          func_ssh_build  && \
          docker build --no-cache --rm -f docker/rpm -t stns:libnss_stns . && \
          docker run -v "$(pwd)"/binary:/go/src/github.com/STNS/libnss_stns/binary -t stns:libnss_stns
            ;;
        d)
          func_ssh_build && \
          docker build --no-cache --rm -f docker/deb -t stns:libnss_stns . && \
          docker run -v "$(pwd)"/binary:/go/src/github.com/STNS/libnss_stns/binary -t stns:libnss_stns
            ;;
    esac
done
shift $((OPTIND - 1))
