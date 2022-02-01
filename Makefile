REVISION := $(shell git rev-parse --short HEAD)

default:
	go build -o bin/kafu -ldflags "-X 'main.revision=$(REVISION)'" \
		main.go

build-tag-push: build tag push

build:
	docker build -t mo7ka/kafu:$(REVISION) .

tag:
	docker tag mo7ka/kafu:$(REVISION) kix.ocir.io/nro0opflyunj/kafu/kafu:$(REVISION)

push:
	docker push kix.ocir.io/nro0opflyunj/kafu/kafu:$(REVISION)
