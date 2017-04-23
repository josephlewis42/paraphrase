# Copyright 2017 Joseph Lewis III <joseph@josephlewis.net>
# Licensed under the MIT License
#
# Builds paraphrase

# Run git tag x.y.z to set up the new version
VERSION=`git describe --tags`
BUILD=`date +%FT%T%z`
BRANCH=`git log --pretty=format:'%h' -n 1` # short branch format

GO=go
GOFLAGS=-ldflags "-w -s -X main.Version=$(VERSION) -X main.Build=$(BUILD) -X main.Branch=$(BRANCH)"
OUTDIR=bin
WIN=env GOOS=windows GOARCH=amd64
LIN=env GOOS=linux GOARCH=amd64

all: paraphrase

paraphrase: outputdir
	$(GO) build $(GOFLAGS) -o $(OUTDIR)/paraphrase main.go
	$(WIN) $(GO) build $(GOFLAGS) -o $(OUTDIR)/paraphrase.exe main.go

version:
	echo "Version: " $(VERSION)
	echo "Build: " $(BUILD)
	echo "Branch: " $(BRANCH)

clean:
	rm -rf $(OUTDIR)

outputdir:
	mkdir -p $(OUTDIR)
	cp LICENSE $(OUTDIR)
	cp README.md $(OUTDIR)
