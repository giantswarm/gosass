all: build

build:
	go build .

install: build
	cp -f gosass /usr/local/bin/gosass

deps:
	go get -u gopkg.in/fsnotify.v1

unittests:
	mkdir -p integration/out
	go test github.com/dailymuse/gosass/compiler
