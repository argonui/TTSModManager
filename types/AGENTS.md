# types

The shared vocabulary for JSON data. Deliberately tiny and dependency-free — every other package
imports it, so anything added here is visible everywhere and nothing here may import from
elsewhere in the project.

## What's in it

- `J` — a JSON object, `map[string]interface{}`. Used throughout in preference to writing the map
  type out, both for brevity and so the intent ("this is decoded JSON") is explicit.
- `ObjArray` — a JSON array of objects, `[]map[string]interface{}`.
- `ConvertToObjArray(v)` — coerces an `interface{}` from `encoding/json` into `[]map[string]interface{}`.

## Why `ConvertToObjArray` exists

`encoding/json` decodes arrays into `[]interface{}`, not into a slice of maps, so every place that
wants to walk an array of objects needs the same type assertion and loop. This function is that
loop, done once.

It accepts an already-converted `[]map[string]interface{}` unchanged, so callers can pass values
that came from either a fresh unmarshal or an earlier conversion without checking first.

**Nil elements are skipped rather than rejected.** Real TTS savegames occasionally contain nulls in
object arrays, and failing the whole build over one is worse than dropping it.

## Testing

No test file. The conversion is exercised heavily through `mod` and `objects`, both of which call it
on every savegame they process. If you add anything non-trivial here, add a `types_test.go` — the
package has no dependencies, so tests need no setup at all.
