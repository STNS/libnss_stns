GO ?= go

all: build

build: #test
	$(GO) build -buildmode=c-shared -o libnss_etcd.so main.go

#test: getdeps
	#$(GO) test $(TESTTARGET)

