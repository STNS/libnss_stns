eval $(docker-machine env dev)
rm -rf build/ubuntu/*
rm -rf releases/libnss_stns*deb
cd ../ssh_stns_wrapper &&
bin/ubuntu_build.sh &&
cd ../libnss_stns &&
cp -p ../ssh_stns_wrapper/build/ubuntu/ssh_stns_wrapper build/ubuntu/ssh-stns-wrapper &&
docker build -f UbuntuDockerfile -t ubuntu:stns . && \
docker run -v "$(pwd)"/releases:/go/src/github.com/pyama86/libnss_stns/releases ubuntu:stns
