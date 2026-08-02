# mod

Orchestrates a whole savegame in both directions. This package owns the *top level* of a mod —
the keys that sit beside `ObjectStates` — and delegates the object tree itself to
[`objects`](../objects/).

## The two halves

| File | Type | Direction |
|---|---|---|
| [generate.go](generate.go) | `Mod` | directory tree → savegame JSON |
| [reverse.go](reverse.go) | `Reverser` | savegame JSON → directory tree |

They are mirror images, and **every change to one needs the matching change to the other.** A
field that `Reverser.Write` externalizes but `Mod.generate` never re-inlines will silently vanish
from rebuilt mods. This symmetry is the single most important property to preserve here.

## The key tables

`generate.go` opens with three lists that drive both directions:

- `ExpectedStr` — top-level keys holding strings (`SaveName`, `Note`, `LuaScript`, …)
- `ExpectedObj` — keys holding a JSON object (`Grid`, `Lighting`, `MusicPlayer`, …)
- `ExpectedObjArr` — keys holding an array of JSON objects (`SnapPoints`, `Decals`, …)

**Adding support for a new savegame field is usually just adding it to the right list.** The
externalize-and-inline machinery is generic over these tables; special-casing a key in code should
be a last resort, and today only `LuaScript`, `XmlUI`, and `SnapPoints` need it.

## Both directions in outline

**Reverse** (`Reverser.Write`) walks the three key tables, and for each key whose value exceeds the
size threshold, writes the value to its own file and replaces the key with `Key_path`. Lua and XML
are handed to [`handler`](../handler/) instead, because they also need bundling. `ObjectStates` is
handed to `objects.Printer`, which returns the object order — stored as `ObjectStates_order`.
Whatever remains is written as `config.json`.

**Generate** (`Mod.GenerateFromConfig`) reads `config.json` and does the inverse: every `Key_path`
is read back from disk and re-inlined as `Key`, Lua and XML go back through `handler` to be
re-bundled, and `objects.ParseAllObjectStates` rebuilds the object array in the order recorded by
`ObjectStates_order`.

## Date and EpochTime

Reversing deletes `Date` and `EpochTime`; generating fills them with the current time. A savegame's
original timestamp says nothing about the mod's *content*, and keeping it would make every rebuild
produce a spurious diff. Any other field that changes without the mod changing should get the same
treatment.

## Object-only mode

When `OnlyObjState` / `OnlyObjStates` is set, both types short-circuit to a single-object path
(`writeOnlyObjStates` / `generateOnlyObjStates`) that skips all top-level savegame keys. This backs
the `--objin`/`--objout` flags for TTS downloadable content, where the file is one object rather
than a savegame.

## `tryPut`

The helper that moves a value from `Key_path` to `Key` during generation. It takes the reader
function as a parameter so the same logic serves strings, objects, and object arrays. Note it
deletes the `_path` key after resolving it — the two forms are never both present in output.

## Testing

`generate_test.go` and `reverse_test.go` drive both types against `tests.FakeFiles`, the in-memory
filesystem from [`tests`](../tests/). Construct a `Mod` or `Reverser` with a fake in every slot;
never touch the real disk from a test in this package.

The round-trip property — reverse then generate reproduces the input — is tested one level up in
[tests/e2e_test.go](../tests/e2e_test.go), because it needs both halves plus `objects`. When you
add a top-level key, add a fixture there too.
