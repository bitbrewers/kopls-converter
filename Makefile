test:
	go test -race -cover ./...

build:
	go build -o builds/kopls-converter

update-deps:
	dep ensure -update

docker-build:
	docker build -t bitbrewers/kopls-converter .

start-dev:
	docker-compose build
	docker-compose up

.PHONY: test docker-run
