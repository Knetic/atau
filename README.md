atau
====

[![Build Status](https://travis-ci.org/Knetic/atau.svg?branch=master)](https://travis-ci.org/Knetic/atau)

Generates code to build executables that can implement RESTful API clients in many languages. In more practical terms, this project aims to be an open-source equivalent to the tools apparently used at Google to build their clientside APIs (see the [discovery libraries](https://developers.google.com/discovery/libraries)).

Why use it?
====

Normally when you build a webservice, you want to implement client libraries for it; redistributables that enable developers to interact with your service easily. Usually, these libraries are only shipped for a language or two (if at all), and are expensive to create and maintain - involving effort from developers who are familiar with every language you want to support. Inevitably, you will never ship a library for every language that a developer may wish to use when interacting with your product. You'll also end up supporting your libraries (including fixing bugs) far into the future, and need to maintain every one of them to match the changes you make to your product.

`atau` seeks to solve that problem. You write a straightforward api description file (no code required) that describes how your webservice functions, what endpoints are available, and what headers/parameters/paths are necessary to interoperate with it. Then, `atau` can generate client libraries for this service in many languages - automatically, with no investment by you. Since the code is automatically generated, it is repeatable, can be tested by the `atau` authors instead of by your dev team, and conforms to whatever best practices the target languages advise.

It's also possible to use `atau` with services that you didn't write. Since you can generate code for any RESTful service, you can transcribe their API to a json file and use `atau` on that. For an example of this, seeing the "lendingclub.json" sample, which describes a few common operations of the LendingClub API. Note that this sample was not provided by LendingClub, it's just a way to use `atau` to generate clients for sites that do not yet publish an API schema document.

How do I use it?
====

Install one of the packages provided in the releases page, or just use the executable built by this repo. Requires an API document in the Discovery format (see the samples, or check out the [docs](https://developers.google.com/discovery/v1/reference/apis#methods)).

In which languages can code be generated?
====

* golang (any)
* C# (.NET 3.5+)
* Python (2+)
* Ruby (1.9+)

Some languages explicitly do not have support from `atau`, detailed below.

Additions
====

For the most part, this library aims to implement the google discovery format. However, there are some extensions and places where this differs from what appears to be their spec.

#### Headers

`atau` allows the specification of headers similar to any other parameter. Every method can specify a "headers" block, which is a map of strings to schemas.

#### Automatic querystring detection

In the Google Discovery API, a parameter given to a method must specify a `location` field, which dictates whether the parameter will be used on the path of the request or in the querystring. `atau` obsoletes this distinction and automatically figures out which parameters are path parameters, and which are query.

Why is my favorite language not listed?
====

* JavaScript: Frontend JS is not on any roadmap, and probably won't ever be - due to the huge amount of libraries and possible ways to implement an API client in frontend JS. However, IO.js and Node.js support seems more unified, and therefore possible.

* Powershell 2+: Support for this seems extremely likely, but will take some time since PS is not yet implemented in [presilo](https://github.com/Knetic/presilo), which powers the schema parsing and code generation for `atau`.

* Java: Unfortunately, `atau` assumes the serialization format of all communications between client and server is JSON. Java does not have a built-in JSON library. Maybe more importantly, even if it did, older and most-prevalent versions of Java (such as 6) do not have such support. The only option would be to have `atau` force the use of a specific JSON library, or to complicate the project by having subsets of java code generation which support different JSON libraries. It would fast become a nightmare to support such things.

* Lua: While the authors have a deep fondness for lua, its usage is rarely as a general-purpose language. Lua does not have first-class support for JSON parsing, nor even for HTTP requests. For all the same reasons that Java support is not planned, lua is also not planned.

* C/C++: These languages are too low-level to have a built-in JSON or HTTP library. The author, as noted above, is not interested in bloating codegeneration on a per-library basis.

Other languages are simply not familiar enough to the author to warrant an implementation of them. Pull requests and iterations on them are welcome, but will first need to be included into [presilo](https://github.com/Knetic/presilo), which powers the schema parsing and code generation for `atau`.

Branching
====

I use green masters, and heavily develop with private feature branches. Full releases are pinned and unchangeable, representing the best available version with the best documentation and test coverage. Master branch, however, should always have all tests pass and implementations considered "working", even if it's just a first pass. Master should never panic.

Activity
====

If this repository hasn't been updated in a while, it's probably because I don't have any outstanding issues to work on - it's not because I've abandoned the project. If you have questions, issues, or patches; I'm completely open to pull requests, issues opened on github, or emails from out of the blue.

Affiliation
====

Neither the author nor this project are associated with Google. This project consumes an open standard of document, but doesn't claim any ownership of the standard. A sample is included for LendingClub API, but the author has no affiliation there either - the sample is just a sample of a public API.
