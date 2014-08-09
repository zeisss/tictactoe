SOURCE := $(shell find . -name '*.go')

tictactoed: $(SOURCE)
	GOPATH=$(GOPATH) go build .