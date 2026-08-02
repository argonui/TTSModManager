# tests

Two things live here: the in-memory filesystem every other package's tests are built on, and the
end-to-end round-trip test.

## `FakeFiles` — [fakefiles.go](fakefiles.go)

Note this is a plain `.go` file, not `_test.go`, and that is deliberate: `mod`, `objects`, and
`bundler` all import it. Keep it importable.

`FakeFiles` implements **every** interface from [`file`](../file/) at once — `TextReader`,
`TextWriter`, `JSONReader`, `JSONWriter`, `DirCreator`, `DirExplorer`. One instance can therefore
fill every slot of a `mod.Mod` or `mod.Reverser`, or several instances can be used to keep
directories distinct (as the e2e test does, with separate fakes for modsettings, objects, and
final output).

State is two maps:

- `Fs` — text files, filename → content
- `Data` — JSON files, filename → `types.J`

**Use `NewFF()` rather than a bare `&FakeFiles{}`** unless you are initializing both maps yourself;
writes into a nil map will panic.

There are no real directories. `ListFilesAndFolders` synthesizes the tree by splitting the `/` in
the keys of `Data`, so a file written to `objects/bag/card.json` implies a folder `objects/bag`
without anything having created it. Tests that care about directory structure should set up `Data`
keys accordingly.

Because the fake is what most tests run against, **its semantics are the semantics under test.**
Where it differs from the real `file` implementations, tests will agree with the fake and disagree
with production. Keeping the two aligned is worth more than any individual test here.

## `e2e_test.go`

`TestAllReverseThenBuild` globs `testdata/e2e/*.json`, and for each savegame: reverses it into
fakes, generates it back out, and compares the result against the original. This is the round-trip
property the whole tool rests on.

To add coverage, drop a savegame JSON into `testdata/e2e/` — the test picks it up with no other
changes. Small, targeted fixtures that isolate one feature are more useful than large real mods.

What the comparison covers:

- `LuaScript` is compared by its set of bundled module names rather than byte-for-byte, since
  bundling is order- and formatting-sensitive.
- `Date` and `EpochTime` are excluded — they are regenerated on every build by design.
- Any `float64` value is excluded. Numeric behavior is asserted in
  [objects/numbersmoother_test.go](../objects/numbersmoother_test.go) instead.
- `bundled_core` and `bundled_lua` are on a deny list; they depend on `src/` files the harness
  doesn't provide.

## Fixture locations

`tests/testdata/e2e/` is the only fixture directory any test reads. The `testdata/` tree at the
repository root is not referenced by any Go code — it exists for manual runs of the binary.
