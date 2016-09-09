NAME := dockertags
SRCS := $(shell find . -name '*.go' -type f)
LDFLAGS := -ldflags="-s -w"

.DEFAULT_GOAL := bin/$(NAME)

bin/$(NAME): $(SRCS)
	go build $(LDFLAGS) -o bin/$(NAME)

.PHONY: clean
clean:
	rm -rf bin/*
