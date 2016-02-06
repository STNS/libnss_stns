FROM golang:latest
ADD . /go/src/github.com/STNS/libnss_stns
RUN go get github.com/tools/godep
WORKDIR /go/src/github.com/STNS/libnss_stns
RUN godep restore
CMD go build libnss_stns.go resource.go passwd.go shadow.go group.go
