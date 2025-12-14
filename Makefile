.PHONY: all build clean test install

# Build variables
BINARY_INSTALLER=bin/seatunnel-installer
BINARY_AGENT=bin/seatunnel-agent
BINARY_CONTROL_PLANE=bin/seatunnel-control-plane

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

all: clean build

build: build-installer build-agent build-control-plane

build-installer:
	@echo "Building installer..."
	@mkdir -p bin
	$(GOBUILD) -o $(BINARY_INSTALLER) ./cmd/installer

build-agent:
	@echo "Building agent..."
	@mkdir -p bin
	$(GOBUILD) -o $(BINARY_AGENT) ./cmd/agent

build-control-plane:
	@echo "Building control plane..."
	@mkdir -p bin
	$(GOBUILD) -o $(BINARY_CONTROL_PLANE) ./cmd/control-plane

clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf bin/

test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

install: build
	@echo "Installing binaries..."
	cp $(BINARY_INSTALLER) /usr/local/bin/
	cp $(BINARY_AGENT) /usr/local/bin/
	cp $(BINARY_CONTROL_PLANE) /usr/local/bin/

run-installer:
	$(GOBUILD) -o $(BINARY_INSTALLER) ./cmd/installer
	$(BINARY_INSTALLER) server

run-agent:
	$(GOBUILD) -o $(BINARY_AGENT) ./cmd/agent
	$(BINARY_AGENT) start

run-control-plane:
	$(GOBUILD) -o $(BINARY_CONTROL_PLANE) ./cmd/control-plane
	$(BINARY_CONTROL_PLANE) server
