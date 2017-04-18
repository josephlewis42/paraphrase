# Paraphrase

Paraphrase is open-source text similarity measurement software inspired by MOSS.
It can compare sets of documents to find content that's been duplicated.
Potential use cases are:

* Checking a collection of student assignments for plagiarism
* Looking for duplicated code in a large repository
* Ensuring you correctly quoted all your sources

## Usage

Paraphrase has a slew of commands, generally your needs will be as simple as
upload then report.

### Uploading

You can upload documents one at a time or in batches

```
$ paraphrase add main.go
Adding: main.go
main.go got id 3

$ paraphrase add examples/*.java
Adding: examples/DocA.java
examples/DocA.java got id 4
Adding: examples/DocB.java
examples/DocB.java got id 5
```

### Reporting

You can run a report for a specific document.
The output is formatted `<document id>: <num matches> (<pct matching>)`.

```
$ paraphrase report 5
4: 21 matches (20%)
2: 101 matches (100%)
...
```

### Testing

You can test various parts of the `paraphrase` tool if you want to see what
it's doing internally.

    xadd        (read only, debug) Dry run of an add.
    xhash       (read only, debug) Print the hashes for a document
    xnorm       (read only, debug) Normalizes files like before they're processed
    xsim        (read only, debug) Calculates the similarity of two documents
    xwinnow     (read only, debug) Print the winnowed hashes

## How does it work?

Check out "[Winnowing: local algorithms for document fingerprinting](https://doi.org/10.1145/872757.872770)"
for information about how paraphrase works behind the scenes.


## Environment

Paraphrase is built to be a fully copy/pasteable binary that can be compiled
for Windows, OSX and Linux.

The only utilities you need to interact with it are a CLI and optionally a
web-browser.

## Algorithms Used

The following algorithms are used/will be used:

* Suffix trees
* [FNV Hashes](https://en.wikipedia.org/wiki/Fowler%E2%80%93Noll%E2%80%93Vo_hash_function)
* Inverted Indexes
* [Winnowing Document Fingerprinting](https://doi.org/10.1145/872757.872770) DOI: 10.1145/872757.872770


## License

Copyright 2017 Joseph Lewis III <joseph@josephlewis.net>

Licensed under the MIT license. See LICENSE for the full text.
