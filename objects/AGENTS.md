# objects

The recursive half of the tool. Where [`mod`](../mod/) handles the flat top level of a savegame,
this package handles `ObjectStates` — an arbitrarily deep tree of objects, each of which may hold
contained objects (a bag of cards) and alternate states (a flippable tile).

## `objConfig` is the tree node

Everything centres on `objConfig`, which holds one TTS object: its raw JSON (`data`), its GUID, and
its children (`subObj`, `states`). It has four methods, one per direction-and-medium:

| Method | From | To |
|---|---|---|
| `parseFromJSON` | savegame JSON | tree |
| `parseFromFile` | disk | tree |
| `print` | tree | savegame JSON |
| `printToFile` | tree | disk |

The parse pair and the print pair must stay symmetric, the same way `mod`'s two halves must.

## File naming

`getAGoodFileName()` produces `<sanitized nickname>.<guid>` — e.g. `AdversaryBag.a1b2c3` — falling
back to the bare GUID when an object has no usable `Nickname` or `Name`. The sanitizing regex keeps
letters and numbers from any language plus `_`, `-`, and `!`, so non-English mods stay readable
rather than collapsing to GUID soup.

The name serves two purposes at once: it is the filename on disk *and* the identity used in
`ObjectStates_order` and `ContainedObjects_order`. It therefore has to be **unique among its
siblings** and **stable across runs** — an unstable name would rewrite every ordering array on
every build. Any change to this function reshuffles filenames across every downstream mod repo.

## Children

Two kinds, both externalized into a subdirectory named after the parent:

- **Contained objects** — `ContainedObjects` in the savegame becomes a `ContainedObjects_path`
  (the subdirectory) plus `ContainedObjects_order` (the array order, which a directory can't record).
- **States** — `States` becomes `States_path`, a map from TTS state ID to filename.

Both are read back by `parseFromFile`, which recurses into the subdirectory.

## Number smoothing

[numbersmoother.go](numbersmoother.go) rounds numeric fields as objects are parsed. TTS writes
positions and rotations at full float precision, so simply nudging a card in-game produces a diff
on the order of `2.9999998807907104` vs `3.0`. Rounding makes diffs mean something.

| Field group | Precision |
|---|---|
| `posX/Y/Z` | 3 decimal places |
| `rotX/Y/Z` | whole degrees, normalized into `[0, 360)` |
| `scaleX/Y/Z` | 2 decimal places |
| color `r/g/b/a` | 5 decimal places |

`smoothArbitrary` handles `x/y/z` triples (snap points, `AltLookAngle`) and deliberately **errors on
any key it wasn't expecting**, rather than passing it through. That is a tripwire: if TTS adds a
field to snap points, the tool fails loudly instead of silently dropping it.

## `db` and the entry points

- `ParseAllObjectStates(...)` — disk → the object array, ordered by the caller-supplied order list.
  It builds a `db` of root objects, then prints them in that order.
- `Printer.PrintObjectStates(...)` — the object array → disk, returning the order to record.

`db.print` requires the order list and the objects on disk to correspond exactly. A mismatch means
someone added an object file without adding it to `ObjectStates_order`, or vice versa.

## Testing

- `objects_test.go` — parse/print of individual objects
- `functional_test.go` — parse-then-print round trips through `tests.FakeFiles`
- `numbersmoother_test.go` — rounding rules in isolation

Number-rounding behavior is asserted *here*, not in the end-to-end test. If you change a precision
value, this package's tests are what will catch the consequences.
