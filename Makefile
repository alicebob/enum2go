.PHONY: all build test vet lint install example

all: build test vet lint

build:
	go build

test:
	go test

vet:
	go vet

lint:
	golint

install:
	go install

example:
	cpp -I ~/src/libnetfilter_conntrack/include/ -I ~/src/libnfnetlink/include/ ~/src/libnetfilter_conntrack/include/libnetfilter_conntrack/libnetfilter_conntrack.h | ./enum2go 

.PHONY: conn.h
conn.h:
	cpp -I ~/src/libnetfilter_conntrack/include/ -I ~/src/libnfnetlink/include/ ~/src/libnetfilter_conntrack/include/libnetfilter_conntrack/libnetfilter_conntrack.h > conn.h

.PHONY: go
go: build
	# ./enum2go nf_conntrack_attr nf_conntrack_attr_grp < conn.h
	./enum2go nf_conntrack_attr_grp < conn.h > gen/hdr.go
	cd gen && go build

.PHONY: demo
demo: build
	# ./enum2go nf_conntrack_attr nf_conntrack_attr_grp < conn.h
	./enum2go -package demo nf_conntrack_attr_grp < conn.h
