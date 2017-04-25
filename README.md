# Paraphrase

Paraphrase is open-source text similarity measurement software inspired by MOSS.
It can compare sets of documents to find content that's been duplicated.
Potential use cases are:

* Checking a collection of student assignments for plagiarism
* Looking for duplicated code in a large repository
* Ensuring you correctly quoted all your sources

## Usage

Paraphrase has a slew of commands, generally your needs will be as simple as
add then report.

### Installing

If you have go installed and your environment variables set up you can run:

```
$ go install github.com/josephlewis42/paraphrase
```

Alternatively, you can download the source and run `make` to build the project.
You'll need go 1.7.

```
$ git clone https://github.com/josephlewis42/paraphrase
<snip>
$ go version
go version go1.7.4 linux/amd64
$ cd paraphrase
$ make
<snip>
$ cd bin
$ ./paraphrase version
Version: 0.0.1-4-ga445927
Build: 2017-04-24T19:34:28-0700
Branch: a445927
```

### Adding documents

You add documents directly, piped from another command, or via `git`.
Each document in the database keeps its name, document text
(comparessed with snappy), internal document id and "path".
The paths aren't valid filesystem paths, but can be used to filter results.

The simplest is a single document:

```
$ paraphrase add ../paraphrase/*.go
Adding: ../paraphrase/db.go
../paraphrase/db.go got id 294 (internal path: ../paraphrase/db.go)
Adding: ../paraphrase/git.go
../paraphrase/git.go got id 295 (internal path: ../paraphrase/git.go)
Adding: ../paraphrase/processing.go
<snip>
```

You can specify a prefix for future searches.
In this case we'll add the version to the front of the path in case we want
to index multiple versions.

```
$ paraphrase add --prefix paraphrasev0.0.1 ../paraphrase/*.go
Adding: ../paraphrase/db.go
../paraphrase/db.go got id 298 (internal path: paraphrasev0.0.1/../paraphrase/db.go)
Adding: ../paraphrase/git.go
../paraphrase/git.go got id 299 (internal path: paraphrasev0.0.1/../paraphrase/git.go)
<snip>
```

Paraphrase works well with other applications.
It can read file paths from stdin (use a `-` to denote stdin).

```
$ find *.java | paraphrase add -
Adding: Main.java
Main.java got id 3 (internal path: Main.java)
<snip>
```

Finally, you can pull directly from a `git` repository.
This isn't recommended for large repos, but will work for smaller ones.
It's a good idea to specify a prefix in this case.
You can use a glob match to specify which files to include from your repo.
In this case, we import only paraphrase's go files.
It will log the commit number and files that are matched to stdout.

```
$ paraphrase addgit --prefix demo --match "**.go" "https://github.com/josephlewis42/paraphrase"
git clone https://github.com/josephlewis42/paraphrase /tmp/paraphrasegit699390947 --recursive
commit a445927d59ecedb22123f6ce8b9c0e345d072e33
Author: Joseph <joseph@josephlewis.net>
Date:   Mon Apr 24 18:33:33 2017 -0700

    cleaned up some commands and started working on git support
<snip>
```

### Viewing Documents

You can list the documents you have in your database.
It will output a <docid> <path> pair separated by a tab.
The `n` flag will limit the number of documents returned

```
$ paraphrase list -n 4
   1	examples/BSD
  10	examples/MIT
<snip>
```

You can also use path matching here to limit your results:

```
$ paraphrase list --match *.go -n 2
 294	../paraphrase/db.go
 295	../paraphrase/git.go

```

You can get information about a document given the id in JSON format.
The hashes are the winnowed hashes that allow matching documents to each other.

```
$ paraphrase get 4
{
    "Id": 4,
    "IndexDate": "2017-04-23T22:50:56-07:00",
    "Path": "examples/MIT",
    "Name": "MIT",
    "Hashes": [
        7143130077406907449,
        645742780794022614,
        649423357929432066,
<snip>
```

You can get the original text of the document:

```
$ paraphrase doctext 10
Permission is hereby granted, free of charge, to any person obtaining a copy
<snip>
```


### Reporting

There are three major ways to get reports out of paraphrase.

1. Search by text (like a search engine)
2. Find similar documents
3. Global report (work in progress)

**Search by text**

You can search for a given string. Note that the results may not be complete
due to the way paraphrase does search matching; it's kind of like a search
engine but not quite.

```
$ paraphrase searchtext "public static void main(String[]"
Search 'public static void main(String[]' was turned into hashes
	5636575385543409200
	1814657868599676395
	4071359337843359789
	4373333261552009081
<snip>
Results:
Id:  150 Matches:   1 Rank: NaN Path: /src/bmod/gui/GuiExtensionPoints.java
Id:  259 Matches:   1 Rank: NaN Path: /src/bmod/plugin/loader/GenericGuiPluginLoader.java
Id:  285 Matches:   1 Rank: NaN Path: /src/edu/du/cs/smartgrid/Common.java
Id:   85 Matches:   1 Rank: NaN Path: /src/Bmod.java
Id:   87 Matches:   1 Rank: NaN Path: /src/bmod/Constants.java
Id:  251 Matches:   1 Rank: NaN Path: /src/bmod/plugin/generic/headless/SmartGridProvider.java
```

**Find similar documents**

This is where things start to get fun. Let's say you have four documents in your
database. The MIT and BSD licenses, and two Hello World programs that have those
licenses.

```
$ paraphrase list
   1	examples/BSD
   2	examples/HelloBsd.java
   3	examples/HelloMit.java
   4	examples/MIT
```

We can search document 2 for similar documents and we expect to see BSD first
because the program is small in comparison to the license.

```
$ paraphrase searchdoc 2
Search results for 2 (examples/HelloBsd.java)
Id:    2 Matches:  74 Rank: 37.75 Path: examples/HelloBsd.java
Id:    1 Matches:  67 Rank: 33.25 Path: examples/BSD
Id:    3 Matches:   6 Rank:  2.75 Path: examples/HelloMit.java
Id:    4 Matches:   1 Rank:  0.25 Path: examples/MIT
```

Things look good there. Notice we have a rank that shows up when documents are
compared. Let's add the BSD license 50 more times and search again limiting
to the top 4 matches.

```
$ seq 50 | xargs -Iz paraphrase add examples/BSD
<snip>
$ paraphrase searchdoc -n 4 2
Search results for 2 (examples/HelloBsd.java)
Id:    2 Matches:  74 Rank:  5.79 Path: examples/HelloBsd.java
Id:    3 Matches:   6 Rank:  2.52 Path: examples/HelloMit.java
Id:   35 Matches:  67 Rank:  1.29 Path: examples/BSD
Id:   39 Matches:  67 Rank:  1.29 Path: examples/BSD
```

Now the other Hello World file comes up higher than the license, even though it has fewer matches.
Paraphrase ranks based on occurrence frequency as well as match count.
This helps eliminate boilerplate code from generating false matches.


**Global report**

The global report is still in progress.
The idea is to generate a report among sets of documents.
For example, to compare one students' homework with the rest of the class:

```
$ paraphrase report compare students/joseph/*.java students/*.java
```



## How does it work?

Check out "[Winnowing: local algorithms for document fingerprinting](https://doi.org/10.1145/872757.872770)"
for information about how paraphrase works behind the scenes.

The following algorithms are used be used:

* TF-IDF (inspired)
* [FNV Hashes](https://en.wikipedia.org/wiki/Fowler%E2%80%93Noll%E2%80%93Vo_hash_function)
* Inverted Indexes
* [Winnowing Document Fingerprinting](https://doi.org/10.1145/872757.872770) DOI: 10.1145/872757.872770


## License

Copyright 2017 Joseph Lewis III <joseph@josephlewis.net>

Licensed under the MIT license. See LICENSE for the full text.
