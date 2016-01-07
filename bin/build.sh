eval $(docker-machine env dev)
rm -rf build/*
rm -rf releases/*
cd ../ssh_stns_wrapper &&
bin/build.sh &&
cd ../libnss_stns &&
cp -p ../ssh_stns_wrapper/bin/build/ssh_stns_wrapper build && \
docker build -t centos:libnss . && \
docker run -v "$(pwd)"/releases:/go/src/github.com/pyama86/libnss_stns/RPM/RPMS centos:libnss
