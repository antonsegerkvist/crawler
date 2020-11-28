GO ?= go

.PHONY: all

all:
	$(GO) run cmd/crawler/main.go