# TTSModManager

A Go CLI that converts a Tabletop Simulator savegame — a single enormous JSON file — into a
directory tree that can live in source control, and converts that tree back into a savegame TTS
can load. Both directions must produce equivalent results; round-tripping is the core promise.

## The guiding principle

**Every file in the output tree should be human-readable.**

A TTS savegame nests megabytes of Lua, XML, and object state inside one JSON blob. The point of
this tool is not merely to split that up — it is to split it up so that a human reading any single
file can see what that file is *for*. Any value large enough to obscure the object holding it gets
moved into its own file, leaving a pointer behind.

That principle is what motivates nearly every design decision below. When adding a feature, ask
"does this keep each file legible on its own?" before asking anything else.

## The `_path` convention

Externalizing a value is expressed in the JSON itself. A key `Foo` holding a large value becomes a
key `Foo_path` holding the filename that now holds it:

```jsonc
// inline — small enough to read in place
"LuaScriptState": "{\"counter\":3}"

// externalized — the value moved to a file
"LuaScriptState_path": "MyDeck.a1b2c3.luascriptstate"
```

Generating reverses this: `Foo_path` is read from disk and re-inlined as `Foo`. Only one of the two
forms is ever present at a time.

Two related conventions record things a directory cannot express on its own:

- `*_order` (`ObjectStates_order`, `ContainedObjects_order`) — arrays are ordered, directory
  listings are not, so the order is stored explicitly.
- `States_path` — a map of TTS state ID to the filename holding that state.

## Size thresholds

Whether a value is externalized is decided by a length check: 80 characters for strings and
scripts, 100 for JSON objects, 200 for object arrays. These numbers are tuning knobs chosen by
feel — they balance "don't bury a big value inside its parent" against "don't explode a mod into
thousands of tiny files." There is no deeper meaning to the specific values.

They are, however, load-bearing for round-trip stability: changing one changes which files exist in
every mod built afterward, producing enormous diffs in downstream repos. Treat a change as a
breaking change.

## The two directions

| Direction | Flag | Entry point | Does |
|---|---|---|---|
| Generate | *(default)* | `mod.Mod.GenerateFromConfig` then `.Print` | directory tree → savegame JSON |
| Reverse | `--reverse` | `mod.Reverser.Write` | savegame JSON → directory tree |

`main.go` is wiring only: it parses flags, constructs the concrete `file.*` implementations, points
them at the right directories, and hands them to `mod`. Logic belongs in packages, not in `main`.

## Mod directory layout

```
<moddir>/
  config.json      the savegame minus everything externalized
  objects/         one .json per object, nested dirs for contained objects
  modsettings/     externalized settings blobs (Grid, Lighting, MusicPlayer, …)
  src/             Lua modules, resolved by require("name") -> src/name.ttslua
  xml/             XML fragments, resolved by <Include src="name"/> -> xml/name.xml
```

## Modes

- **Full mod** — the default; reads/writes `config.json` plus the whole tree.
- **Object only** — `--objin` / `--objout` operate on a single object-state JSON rather than a
  whole savegame, for TTS downloadable content. Both flags must be set or neither.
- **Saved object** — `--savedobj` wraps output in the TTS "saved object" envelope. See
  [saved-object-feature.md](saved-object-feature.md).
- **`--writesrc`** — when reversing, also write the unbundled `require`d modules to `src/`.
- **`--bonusdir`** — an additional root searched when resolving `require`/`Include`, so shared
  libraries can live outside the mod.

## Packages

| Package | Responsibility |
|---|---|
| `main` | flag parsing and dependency wiring only |
| [`mod`](mod/) | orchestrates a whole savegame in both directions |
| [`objects`](objects/) | the recursive object tree; contained objects, states, GUID naming, number smoothing |
| [`handler`](handler/) | shared decision logic for "inline or file?" across Lua and XML |
| [`bundler`](bundler/) | Lua `require` bundling and XML `Include` bundling, both directions |
| [`file`](file/) | every filesystem operation, and the interfaces the rest of the code depends on |
| [`types`](types/) | `J` (a JSON object) and array conversion helpers |
| [`tests`](tests/) | in-memory fake filesystem and the end-to-end round-trip test |

## Commands

```bash
go build ./... && go test ./... && go vet ./...
```

Go 1.18 is the floor — set in `go.mod` and pinned in both CI workflows. Newer toolchains build it
fine, but don't raise the `go` directive without also updating `.github/workflows/go.yml` and the
two `.slsa-goreleaser*.yml` release configs.

## Conventions

- **All filesystem access goes through `file/`.** Outside of `file/` and `main.go`, no package
  imports `os` or `io`. `objects/` imports `path` purely for string joining, never for I/O. This is
  what makes the whole pipeline testable against an in-memory fake — preserve it.
- **Depend on interfaces, not implementations.** `mod`, `objects`, and `handler` accept
  `file.TextReader`, `file.JSONWriter`, and friends. Never take a `*file.JSONOps`.
- **`types.J` is the universal JSON object type**, an alias for `map[string]interface{}`.
- **Wrap errors with the call that produced them**: `fmt.Errorf("EncodeFromFile(%s): %v", name, err)`.
  The resulting chains are how failures deep in a nested object tree stay diagnosable.
- **Keys are data.** The lists in [mod/generate.go](mod/generate.go) (`ExpectedStr`, `ExpectedObj`,
  `ExpectedObjArr`) enumerate the top-level savegame keys and how each is treated. Adding support
  for a new savegame field usually means adding it to one of those lists, not writing new logic.

## Before assuming behavior is intentional

This codebase has known open defects tracked as GitHub issues. If something here looks wrong or
surprising, check the issue tracker before treating it as an invariant to preserve.
