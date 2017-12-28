test:
	go test -race -cover ./...

update-deps:
	dep ensure -update

docker-build:
	docker build -t bitbrewers/kopls-converter .

docker-run: docker-build
	docker run --rm -it \
	-p 8000:8000 \
	-v $(shell pwd):/go/src/github.com/bitbrewers/kopls-converter \
	bitbrewers/kopls-converter \
	fresh -c .freshconf

.PHONY: test update-deps docker-build docker-run
