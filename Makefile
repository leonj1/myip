ifndef VERSION
  VERSION=$(shell git rev-parse --short HEAD)
endif

all: docker

build:
	go build -v -o myip myip.go

docker: build
	docker build -t www.dockerhub.us/myip:$(VERSION) .

push:
	docker push www.dockerhub.us/myip:$(VERSION)

