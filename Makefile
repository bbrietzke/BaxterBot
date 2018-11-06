.PHONY: dependencies clean clean-all test

run: dependencies test
	$(GO_RUN) main.go

run2: dependencies test
	$(GO_RUN) main.go --swarm 22000 --http 22080 --join localhost:8080/swarm

run3: dependencies test
	$(GO_RUN) main.go --swarm 22000 --http 22080 --join localhost:8080/swarm

dependencies: Gopkg.toml
	$(DEP) ensure -v

clean:
	$(GO) clean

clean-all: clean
	git clean -xfd

$(BINDIR):
	mkdir $(BINDIR)
	
$(DEP): 
	brew install dep

Gopkg.toml: $(DEP)
	$(DEP) init

swarm: pkg/swarm/swarm.pb.go

pkg/swarm/swarm.pb.go:
	PATH=$(TPATH) $(PROTOC) -I/usr/local/include -I. -I$(GOPATH)/src -I proto/ swarm.proto --go_out=plugins=grpc:pkg/swarm

test:
	@echo "Test something here"

test-template:
	$(GO_TEST) github.com/bbrietzke/BaxterBot/*

build-pi: $(BINDIR) test
	GOOS=linux GOARCH=arm GOARM=5 $(GO_BUILD) -o $(BINDIR)/pi/BaxterBot


BIN := /usr/local/bin
DEP := $(BIN)/dep
PROTOC := $(BIN)/protoc
TPATH := $(PATH):$(GOPATH)/bin

GO := $(BIN)/go
GOBIN := $(GOPATH)/bin
GO_RUN := $(GO) run
GO_TEST := $(GO) test
GO_BUILD := $(GO) build

LOCAL := $(shell pwd)
BINDIR := $(LOCAL)/bin
GIT_HASH := $(shell git show --format="%h" --no-patch)

