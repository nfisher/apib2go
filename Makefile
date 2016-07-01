SHELL := /bin/sh
EXE   := apib2go
SRC   := $(wildcard *.go)
COVER := cover.out

.DEFAULT_GOAL := all

.PHONY: all
all: $(EXE)

# Build the exe.
$(EXE): $(COVER)
	go build -v

# Run vet, test, and display test coverage by function.
$(COVER): $(SRC)
	go vet -v
	go test -v -covermode=count -coverprofile=$(COVER)
	go tool cover -func=$(COVER)

# Runs the html based coverage tool.
.PHONY: cov
cov: $(COVER)
	go tool cover -html=$(COVER)

# Clean up the project files.
.PHONY: clean
clean:
	go clean -v
	$(RM) -f $(COVER)
