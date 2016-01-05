GO ?= go

all: build

build:
	$(GO) get github.com/BurntSushi/toml
	$(GO) build -buildmode=c-shared -o libnss_stns.so
