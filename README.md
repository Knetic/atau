atau
====

Generates code to build executables that can implement RESTful API clients in many languages. In more practical terms, this project aims to be an open-source equivalent to the tools apparently used at Google to build their clientside APIs (see the [discovery libraries](https://developers.google.com/discovery/libraries)).

Why use it?
====

Normally when you build a webservice, you want to implement client libraries for it; redistributables that enable developers with interact with your service easily. Usually, these libraries are only shipped for a language or two (if at all), and are expensive to create and maintain. Inevitably, you will never ship a library for every language that a developer may wish to use when interacting with your product.

`atau` seeks to solve that problem. You write a straightforward api description file (no code required) that describes how your webservice functions, what endpoints are available, and what headers/parameters/paths are necessary to interoperate with it. Then, `atau` can generate client libraries for this service in many languages - automatically, with no investment by you. Since the code is automatically generated, it is repeatable, can be tested by the `atau` authors instead of by your dev team, and conforms to whatever best practices the target languages advise.

In which languages can code be generated?
====

* golang
* C#
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

* JavaScript: Frontend JS is not on any roadmap, and probably won't ever be. However, IO.js and Node.js support seems likely.

* Powershell 2+: Support for this seems extremely likely, but will take some time since PS is not yet implemented in [presilo](https://github.com/Knetic/presilo), which powers the schema parsing and code generation for `atau`.

* Java: Unfortunately, `atau` assumes the serialization format of all communications between client and server is JSON. Java does not have a built-in JSON library. Maybe more importantly, even if it did, older and most-prevalent versions of Java (such as 6) do not have such support. The only option would be to have `atau` force the use of a specific JSON library, or to complicate the project by having subsets of java code generation which support different JSON libraries. It would fast become a nightmare to support such things.

* Lua: While the authors have a deep fondness for lua, its usage is rarely as a general-purpose language. Lua does not have first-class support for JSON parsing, nor even for HTTP requests. For all the same reasons that Java support is not planned, lua is also not planned.

* C/C++: These languages are too low-level to have a built-in JSON or HTTP library. The author, as noted above, is not interested in bloating codegeneration on a per-library basis.

Other, less-used languages are simply not familiar enough to the author to warrant an implementation of them. Pull requests and iterations on them are welcome, but will first need to be included into [presilo](https://github.com/Knetic/presilo), which powers the schema parsing and code generation for `atau`.
