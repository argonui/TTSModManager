# handler

Shared decision logic for scripts. Four situations need the same set of choices made — Lua and XML,
each at the savegame root and on every object — and duplicating them four ways is what this package
exists to prevent.

## The decisions it owns

For a script being written to disk:

1. Unbundle it into its component modules.
2. Is the root module long enough to deserve its own file, or should it stay inline?
3. Which writer gets it — the default one, or the `src/` writer for `require`d modules?

For a script being read back:

1. Is the content inline (`LuaScript`) or in a file (`LuaScript_path`)?
2. Re-bundle it, resolving `require`/`Include` through the reader.

## Configuration, not subclassing

`Handler` holds the parts that differ, and `NewLuaHandler()` / `NewXMLHandler()` fill them in:

| Field | Lua | XML |
|---|---|---|
| `key` / `keypath` | `LuaScript` / `LuaScript_path` | `XmlUI` / `XmlUI_path` |
| `extension` | `.ttslua` | `.xml` |
| `bundle` / `unbundle` | `bundler.Bundle` / `UnbundleAll` | `bundler.BundleXML` / `UnbundleAllXML` |

A third kind of bundled content would be a third constructor, not new logic.

## `HandleAction`

Neither method mutates the caller's map. They return a `HandleAction` describing what the caller
should do — set `Key` to `Value`, or do nothing if `Noop` — and the caller applies it.

This matters because the *right* key depends on the outcome: a script written to a file needs
`LuaScript_path` set, while one kept inline needs `LuaScript`. Returning the decision instead of
performing it keeps that from being duplicated at each of the four call sites.

Callers should delete both `key` and `keypath` before applying the action, so the two forms can
never both be present.

## `SrcWriter`

A nil `SrcWriter` means "don't write the `require`d modules to `src/`" — this is how the
`--writesrc` flag is expressed. The root module is still written via `DefaultWriter` either way;
only the dependencies are skipped.

## Testing

This package has no test file of its own. Its behavior is covered indirectly through `mod` and
`objects`, which exercise all four call sites, and through the end-to-end round trip in
[`tests`](../tests/). Direct unit tests would be a welcome addition — the package is pure logic over
injected readers and writers, so it needs no fixtures beyond a fake.
