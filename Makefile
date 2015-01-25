.PHONY: all build test vet lint install example

all: build test vet

build:
	go build

test:
	go test

vet:
	go vet

lint:
	golint .

install: build
	go install

# some example
LINUXSRC=~/src/linux-3.14.7

.PHONY: p9
p9: build
	# ./enum2go nf_conntrack_attr nf_conntrack_attr_grp < conn.h
	./enum2go -I ${LINUXSRC}/include p9_msg_t < ${LINUXSRC}/include/net/9p/9p.h
