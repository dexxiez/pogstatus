GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GORUN=$(GOCMD) run
BINARY_NAME=pog_status
BINARY_DIR=bin
INSTALL_DIR=/usr/local/bin

all: build

build:
	mkdir -p $(BINARY_DIR)
	$(GOBUILD) -o $(BINARY_DIR)/$(BINARY_NAME) -v

clean:
	$(GOCLEAN)
	rm -rf $(BINARY_DIR)

run:
	$(GORUN) .

build_and_run: build
	./$(BINARY_DIR)/$(BINARY_NAME)

install: build
	cp $(BINARY_DIR)/$(BINARY_NAME) $(INSTALL_DIR)

uninstall:
	rm -f $(INSTALL_DIR)/$(BINARY_NAME)

.PHONY: all build clean run build_and_run install uninstall
