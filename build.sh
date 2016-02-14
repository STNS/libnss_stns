eval $(docker-machine env dev)
docker build --rm -t stns:libnss_stns .
docker run -v "$(pwd)"/binary:/go/src/github.com/STNS/libnss_stns/binary -t stns:libnss_stns
