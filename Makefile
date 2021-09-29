all: build docker
build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

docker:
	sudo docker build -f Dockerfile .

