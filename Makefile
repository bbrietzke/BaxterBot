.PHONY: dependencies clean clean-all test

run: dependencies test
	$(GO_RUN) main.go

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


BIN := /usr/local/bin
DEP := $(BIN)/dep
PROTOC := $(BIN)/protoc
TPATH := $(PATH):$(GOPATH)/bin

GO := $(BIN)/go
GOBIN := $(GOPATH)/bin
GO_RUN := $(GO) run
GO_TEST := $(GO) test

LOCAL := $(shell pwd)
BINDIR := $(LOCAL)/bin
GIT_HASH := $(shell git show --format="%h" --no-patch)

