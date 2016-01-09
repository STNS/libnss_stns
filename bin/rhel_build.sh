eval $(docker-machine env dev)
rm -rf build/rhel/*
rm -rf releases/libnss_stns*rpm
cd ../ssh_stns_wrapper &&
bin/rhel_build.sh &&
cd ../libnss_stns &&
cp -p ../ssh_stns_wrapper/build/rhel/ssh_stns_wrapper build/rhel && \
docker build -f RhelDockerfile -t centos:stns . && \
docker run -v "$(pwd)"/releases:/go/src/github.com/pyama86/libnss_stns/releases centos:stns
