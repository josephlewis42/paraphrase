GO=go
GOFLAGS=
OUTDIR=bin

all: paraphrase

paraphrase: outputdir
	$(GO) $(GOFLAGS) build -o $(OUTDIR)/paraphrase main.go 

clean:
	rm -rf $(OUTDIR)

outputdir:
	mkdir -p $(OUTDIR)
	cp LICENSE $(OUTDIR)
	cp README.md $(OUTDIR)
