all: build

build:
	CGO_ENABLED=0 go build -o dist/ccinf github.com/oxdz/comicinfo/cmd/ccinf/