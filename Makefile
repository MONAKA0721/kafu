REVISION := $(shell git rev-parse --short HEAD)

default:
	go build -o bin/kafu -ldflags "-X 'main.revision=$(REVISION)'" \
		main.go
