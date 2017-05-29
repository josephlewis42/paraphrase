# Copyright 2017 Joseph Lewis III <joseph@josephlewis.net>
# Licensed under the MIT License
#
# Builds paraphrase

# Run git tag x.y.z to set up the new version
VERSION=`git describe --tags`
BUILD=`date +%FT%T%z`
BRANCH=`git log --pretty=format:'%h' -n 1` # short branch format

# colors for BASH
GREEN=\033[1;32m
CYAN=\033[1;36m
NC=\033[0m # No Color

GO=go
GOFLAGS=-ldflags "-w -s -X main.Version=$(VERSION) -X main.Build=$(BUILD) -X main.Branch=$(BRANCH)"
OUTDIR=bin
WIN=env GOOS=windows GOARCH=amd64
LIN=env GOOS=linux GOARCH=amd64
RELEASEOS=linux windows darwin
RELEASEDIR=release

BUILDCMD=$(GO) build $(GOFLAGS) -o $(OUTDIR)/paraphrase main.go

all: paraphrase

paraphrase: dependencies outputdir
	@echo
	@echo "${CYAN}Building paraphrase for current OS ${NC}"
	@echo "${GREEN}use make release for all supported OSs${NC}"
	@echo

	$(BUILDCMD)

$(RELEASEOS): clean dependencies outputdir
	@echo
	@echo "${CYAN}Building release" $@ "${NC}"
	mkdir -p $(RELEASEDIR)
	env GOOS=$@ GOARCH=amd64 $(BUILDCMD)
	zip -r $(RELEASEDIR)/paraphrase-$@-$(VERSION).zip bin/

release: version cleanrelease $(RELEASEOS)

cleanrelease:
	@echo
	@echo "${CYAN}Cleaning release directory ${NC}"
	rm -rf $(RELEASEDIR)

version:
	@echo
	@echo "${CYAN}Paraphrase Version ${NC}"
	@echo "Version: " $(VERSION)
	@echo "Build: " $(BUILD)
	@echo "Branch: " $(BRANCH)
	@echo

clean:
	@echo "${CYAN}Cleaning output directory ${NC}"
	rm -rf $(OUTDIR)

dependencies:
	@echo "${CYAN}Getting go dependencies ${NC}"
	$(GO) get ./...

outputdir:
	mkdir -p $(OUTDIR)
	cp LICENSE $(OUTDIR)
	cp README.md $(OUTDIR)
