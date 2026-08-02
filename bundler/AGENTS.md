# bundler

Turns a tree of source files into the single blob TTS wants, and back. Two independent formats live
here — Lua and XML — sharing only the convention that the entry-point content is stored under the
key `Rootname` in an unbundled map.

TTS accepts exactly one Lua string and one XML string per object. Everything this package does
exists so that authors can instead write many small files and have them assembled on build.

## Lua — [luabundler.go](luabundler.go)

`require("name")` resolves to `name.ttslua`, fetched through the injected `file.TextReader` (which
is usually searching several directories — see [file](../file/)). `Bundle` walks requires
breadth-first, collecting every transitively reachable module, then emits the `luabundle` wrapper:
each module wrapped in `__bundle_register("name", function(...) ... end)`, with
`return __bundle_require("__root")` as the entry point.

**The emitted format must stay wire-compatible with luabundle 1.6.0.** It is the de-facto standard
across the TTS scripting ecosystem, and mods built here are consumed by tools that expect it. The
`metaprefix` constant is a verbatim copy of that runtime — treat it as a fixed external contract,
not as code to tidy up.

`Bundle` short-circuits via `IsBundled` when handed already-bundled Lua, and skips bundling
entirely when a script has no requires — a mod with no dependencies stays plain readable Lua rather
than gaining a wrapper it doesn't need.

`UnbundleAll` reverses this, returning module name → source. `Unbundle` is the convenience form
that returns only the root.

## XML — [xmlbundler.go](xmlbundler.go)

`<Include src="name"/>` resolves to `name.xml` and is replaced inline, recursively, with the
included content indented to match the include tag's own indentation.

The bundled form keeps a matched pair of `<!-- include name -->` comments around each inlined
block. Those markers are what make the operation reversible: `UnbundleAllXML` scans for them,
pulls the enclosed lines back out, un-indents them, and restores the `<Include/>` tag. Removing or
reformatting those comments breaks unbundling.

`../` sequences are stripped from include names when unbundling, so a mod cannot direct writes
outside its own tree.

## `Rootname`

The map key under which the entry-point content is stored, used by both formats. Callers get the
root back with `scripts[bundler.Rootname]` rather than a hardcoded string, and
[`handler`](../handler/) relies on it being present in every unbundled result.

## Testing

`luabundler_test.go` and `xmlbundler_test.go` cover both directions. The strongest tests here are
round trips — bundle then unbundle, or unbundle then bundle, and compare — since that is the
property the rest of the tool depends on.

Two end-to-end fixtures (`bundled_core`, `bundled_lua`) are currently on the deny list in
[tests/e2e_test.go](../tests/e2e_test.go) because they require `src/` files the harness doesn't
supply. Wiring those up would give bundling real integration coverage.
