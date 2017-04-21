GO=go
GOFLAGS=
OUTDIR=bin
WIN=env GOOS=windows GOARCH=amd64
LIN=env GOOS=linux GOARCH=amd64

all: paraphrase

paraphrase: outputdir
	$(GO) $(GOFLAGS) build -o $(OUTDIR)/paraphrase main.go 
	$(WIN) $(GO) $(GOFLAGS) build -o $(OUTDIR)/paraphrase.exe main.go 

clean:
	rm -rf $(OUTDIR)

outputdir:
	mkdir -p $(OUTDIR)
	cp LICENSE $(OUTDIR)
	cp README.md $(OUTDIR)
