default: build-docker

build:
	go build

build-docker: build
	docker build -t yaccio/orchid-server .
