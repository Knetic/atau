atau
====

Generates code to build executables that can implement RESTful API clients in many languages. In more practical terms, this project aims to be an open-source equivalent to the tools apparently used at Google to build their clientside APIs (see the [discovery libraries](https://developers.google.com/discovery/libraries)).

Should I use it?
====

Cool your boots. `atau` is in a very early stage. While its aims are not unattainable, it may take some nights and weekends to get it there. It's not fit for production use just yet.

Additions
====

For the most part, this library aims to implement the google discovery format. However, there are some extensions and places where this differs from what appears to be their spec.

#### Headers

`atau` allows the specification of headers similar to any other parameter. Every method can specify a "headers" block, which is a map of strings to schemas.
