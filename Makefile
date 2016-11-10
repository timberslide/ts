all:
	go build

install:
	go install

release:
	go build
	tar -zcvf ts-macos.tar.gz ts

debian-release:
	docker build -f Dockerfile.debian -t ts-debian .
	docker run --name ts-debian ts-debian
	docker cp ts-debian:/ts ts
	chmod +x ts
	tar -zcvf ts-linux-amd64.tar.gz ts
