ARTIFACT_NAME=tinder-cli

all: build

install_deps:
	go get -d ./...
build: install_deps
	go build -o $(ARTIFACT_NAME) *.go
