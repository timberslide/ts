all:
	go build

install:
	go install

release:
	go build
	tar -zcvf ts.tar.gz ts
