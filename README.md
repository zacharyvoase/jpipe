# jpipe

jpipe is a Go rewrite of my Python [jsonpipe][] tool. It's written to be easier
to distribute (thanks to Go's static linking). It should be faster, too—but
don't rely on this tool for performance!

  [jsonpipe]: https://github.com/zacharyvoase/jsonpipe


## Installation

If you have a `GOPATH` you can just:

    go install github.com/zacharyvoase/jpipe


## Example

A `<pre>` is worth a thousand words. For simple JSON values:

    $ echo '"Hello, World!"' | jpipe
    /   "Hello, World!"
    $ echo 123 | jpipe
    /   123
    $ echo 0.25 | jpipe
    /   0.25
    $ echo null | jpipe
    /   null
    $ echo true | jpipe
    /   true
    $ echo false | jpipe
    /   false

The 'root' of the object tree is represented by a single `/` character,
and for simple values it doesn't get any more complex than the above.
Note that a single tab character separates the path on the left from the
literal value on the right.

Composite data structures use a hierarchical syntax, where individual
keys/indices are children of the path to the containing object:

    $ echo '{"a": 1, "b": 2}' | jpipe
    /   {}
    /a  1
    /b  2
    $ echo '["foo", "bar", "baz"]' | jpipe
    /   []
    /0  "foo"
    /1  "bar"
    /2  "baz"

For an object or array, the right-hand column indicates the datatype,
and will be either `{}` (object) or `[]` (array). For objects, the order
of the keys is preserved in the output.

The path syntax allows arbitrarily complex data structures:

    $ echo '[{"a": [{"b": {"c": ["foo"]}}]}]' | jpipe
    /   []
    /0  {}
    /0/a    []
    /0/a/0  {}
    /0/a/0/b    {}
    /0/a/0/b/c  []
    /0/a/0/b/c/0    "foo"


## Caveat: Path Separators

Because the path components are separated by `/` characters, an object
key like `"abc/def"` would result in ambiguous output. jpipe will
throw an error if this occurs in your input, so that you can recognize
and handle the issue. To mitigate the problem, you can choose a
different path separator:

    $ echo '{"abc/def": 123}' | jpipe -k '☃'
    ☃   {}
    ☃abc/def    123

The Unicode snowman is chosen here because it's unlikely to occur as
part of the key in most JSON objects, but any character or string (e.g.
`:`, `::`, `~`) will do.


## Unlicense

This is free and unencumbered software released into the public domain.

Anyone is free to copy, modify, publish, use, compile, sell, or distribute this
software, either in source code form or as a compiled binary, for any purpose,
commercial or non-commercial, and by any means.

In jurisdictions that recognize copyright laws, the author or authors of this
software dedicate any and all copyright interest in the software to the public
domain. We make this dedication for the benefit of the public at large and to
the detriment of our heirs and successors. We intend this dedication to be an
overt act of relinquishment in perpetuity of all present and future rights to
this software under copyright law.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.  IN NO EVENT SHALL THE
AUTHORS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

For more information, please refer to <http://unlicense.org/>
