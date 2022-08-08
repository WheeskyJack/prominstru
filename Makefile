.PHONY: buildalpine
buildalpine:
	GOOS=linux go build -o main .

.PHONY: build
build:
	go build -o main .

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: vendor
vendor:
	go mod vendor

.PHONY: download
download:
	go mod download

.PHONY: image
image: buildalpine
	docker build -t queryexport:1.0.0 -f Dockerfile.alpine .