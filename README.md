# Paraphrase

**Paraphrase** is open-source text similarity measurement software inspired by MOSS.
It can compare sets of documents to find content that's been duplicated.
Potential use cases are:

* Checking a collection of student assignments for plagiarism
* Looking for duplicated code in a large repository
* Ensuring you correctly quoted all your sources

## Usage

**Paraphrase** has a slew of commands, generally your needs will be as simple as
upload then report.

### Uploading

You can upload documents one at a time or in batches using pipes:

```
$ paraphrase upload path/to/my/file.ext
path/to/my/file.ext
Document uploaded with id 125599466215815689

$ find -name "*.java" | paraphrase upload
./Hello.java
Document uploaded with id 125599466215815689
./World.java
Document uploaded with id 5439377427244003676
```

### Reporting

You can run a report for a specific document.
The output is formatted `<document id>: <num matches> (<pct matching>)`.

```
$ paraphrase report 5439377427244003676
5439377427244003676: 21 matches (37%)
1239371856205487029: 7 matches (7%)
...
```

### Testing

You can test various parts of the `paraphrase` tool if you want to see what
it's doing internally.

 normalize		prints the normalized version of the doc
 fingerprint	prints out the fingerprints of the doc
 winnow		prints out the winnowedfingerprints of the doc
 similarity 	prints the overlap between two local files
 upload 		uploads a file to the database
 body 			shows the body of a given file
 report		generates a similarity report for the document with the given id
 hashes		gets the hashes for the document with the given id
 info			prints out info for the document with the given id
 list			lists the ids of all the loaded documents


## How does it work?



## Accuracy

Ultimately, **Paraphrase** is a tool meant to find similarity between documents.
It works very much the same way search engines do.


## Environment

**Paraphrase** is built to be a fully copy/pasteable binary that can be compiled
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

Licensed under the MIT license.
