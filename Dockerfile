FROM golang:latest
ADD . /go/src/github.com/STNS/libnss_stns
WORKDIR /go/src/github.com/STNS/libnss_stns
RUN go get github.com/tools/godep && godep restore
CMD go test ./... && GOARCH=amd64 go build -o binary/libnss-stns.so -buildmode=c-shared main.go resource.go passwd.go shadow.go group.go
