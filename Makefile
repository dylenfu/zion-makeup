SHELL=/bin/bash

# Go parameters
GOCMD=GO111MODULE=on go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test

compile:
	@$(GOBUILD) -o ./build/setup main.go

compile-linux:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o ./build/setup-linux main.go

run:
	cp config/config.json build/config.json
	@echo nodes number $(nodes)
	./build/setup -nodes=$(nodes) -config=build/config.json

clean:
	rm -rf build/*