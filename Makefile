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
	@echo nodes number $(nodes)
	./build/setup -nodes=$(nodes) -config=build/config.json

clean:
	rm -rf build/nodes build/alloc-nodes.json build/extra.dat build/minerlist.txt build/static-nodes.json build/setup