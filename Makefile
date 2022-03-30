SHELL=/bin/bash

# Go parameters
GOCMD=GO111MODULE=on go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
ENV=$(ONROBOT)

compile:
	@$(GOBUILD) -o ./build/$(ENV)/setup main.go

compile-linux:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o ./build/$(ENV)/setup-linux main.go

run:
	@echo nodes number $(nodes)
	./build/$(ENV)/setup -nodes=$(nodes) -env=$(ENV) -config=build/$(ENV)/config.json

clean:
	rm -rf build/$(ENV)/nodes build/$(ENV)/genesis.json build/$(ENV)/alloc-nodes.json build/$(ENV)/extra.dat build/$(ENV)/minerlist.txt build/$(ENV)/static-nodes.json build/$(ENV)/setup build/$(ENV)/minerlist.sh