# Deprecation is read from the go comment

The generator in `pkg/middleware/gen` marks a field of the script API as
`@deprecated` when the go source that declares it says so in a comment. It does
not infer the deprecation from the field names, and it does not keep a list of
deprecated fields in this repository. This document records why, because both
alternatives are the obvious thing to reach for and one of them was written
first.

## Scope

Applies to the typedefs the generator produces for `doc/jsdoc.js` and
`doc/script-env.d.ts` — the script-facing documentation the web-ui completes
from. The fields come from structs this repository does not own: the platform
models and the device-repository model, plus `pkg/model` for the option shapes.

**Not this if** the question is whether this library folds a deprecated alias
into its replacement before forwarding it. It does not, and that is a separate
decision about the criteria themselves; the README section "aspects in a filter
criteria" has it. This document is only about how the generated documentation
learns that a field is deprecated.

**Also not this if** the question is where a deprecation is *decided*. That
happens in the model that declares the field, and this repository only reports
what it finds. A field this repository shows as current is a field whose model
does not mark it, which cannot be fixed here.

## How it works

`pkg/middleware/gen/jsdoc/deprecation.go` resolves the declaring struct of a
field through reflection, finds that package's source directory with
`build.Import(pkgPath, ".", build.FindOnly)` — for a dependency that is the
module cache — and parses it with `parser.ParseDir`. A field counts as
deprecated when a comment line above it or beside it opens with the word
`deprecated`. The text after that word, minus an introducing `:`, `-` or `,`,
becomes the note.

```go
AspectId  string   `json:"aspect_id"` //deprecated: alias for a single element AspectIds
AspectIds []string `json:"aspect_ids,omitempty"`
```

becomes

```ts
/** @deprecated alias for a single element AspectIds */
aspect_id: string;
aspect_ids: string[];
```

The marker is the word, not the colon: a field with nothing further to say is
usually marked with a bare `deprecated`, and that has to count. The word has to
end there, so a comment about a `deprecatedName` does not mark the field it sits
on. `deprecationOfLine` is the whole rule, and
`pkg/middleware/gen/jsdoc/deprecated_test.go` pins each spelling, including the
two that must not match.

Reading comments of a *dependency* needs its source, so this step depends on the
module cache being populated. It is, whenever the package compiles, which is
always the case when the generator runs.

## Rejected: derive it from the field names

The first version marked a field `<x>` as deprecated when the same struct also
carried `<x>s` of the matching list type — `aspect_id` next to `aspect_ids`,
`aspect_node` next to `aspect_nodes`. It produces exactly the right answer for
every field in the current models, which is what makes it tempting.

It is still wrong, in both directions. A pair like `path` and `paths` of the
matching type would be marked without being an alias, and a deprecated field
with no replacement named after it — the common case once a field is simply
going away — would never be marked at all. The rule encodes a naming accident,
and the accident holding today says nothing about the next field.

## Rejected: a list kept in this repository

An explicit list of deprecated fields would be exact and easy to read. It is a
second place for a fact that already exists in the model, and the two go out of
sync in one direction only: the model deprecates a new field, nobody remembers
this list, and the generated documentation keeps presenting the field as
current. Nothing fails, so nothing reports it.

The comment is in the same file as the field and changes in the same commit as
the deprecation. That is the property the list cannot have.

## What this costs

The generator parses the source of every package it reflects over, once per run
and cached in memory. It also means a malformed or unreachable package source
yields no deprecations rather than an error — the generated documentation is
then merely less informative, which is the direction to fail in for a
documentation generator.
