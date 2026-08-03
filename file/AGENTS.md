# file

The only package that touches the filesystem, and the home of the interfaces every other package
depends on. This separation is what makes the whole pipeline testable: `mod`, `objects`, and
`handler` never import `os`, so tests can hand them an in-memory fake instead
([`tests.FakeFiles`](../tests/)).

**Keep `os` and `io` imports inside this package.** If another package needs to read or write
something, add a method to an interface here rather than reaching for `os` there.

## The interfaces

| Interface | Method(s) | Implementation |
|---|---|---|
| `TextReader` | `EncodeFromFile` | `TextOps` |
| `TextWriter` | `EncodeToFile` | `TextOps` |
| `JSONReader` | `ReadObj`, `ReadObjArray` | `JSONOps` |
| `JSONWriter` | `WriteObj`, `WriteObjArray`, `WriteSavedObj` | `JSONOps` |
| `DirCreator` | `CreateDir`, `Clear` | `DirOps` |
| `DirExplorer` | `ListFilesAndFolders` | `DirOps` |

They are deliberately narrow. A consumer that only reads Lua takes a `TextReader`, so it is obvious
from the signature that it cannot write anything.

## Many read paths, one write path

`NewTextOpsMulti(readDirs, writeDir)` is the interesting constructor: reads try each directory in
order and return the first hit, while writes always go to a single directory. This is what lets
`require("shared/util")` resolve against the mod's `src/`, the mod's `objects/`, *and* an external
`--bonusdir`, without the caller knowing which one supplied it. `NewTextOps(dir)` is the
single-directory case.

Read order is significant — earlier directories shadow later ones.

## Trailing newlines

Reads strip one trailing `\n`; writes append one. Files on disk therefore end with a newline, the
way POSIX tools expect, without that newline leaking into the mod content and appearing as a
spurious character in TTS. The two halves are a matched pair — changing one without the other adds
or eats a newline on every round trip.

## `Clear` and its safety checks

`DirOps.Clear()` deletes and recreates a directory. Because it is destructive it runs two guards
before touching anything.

**Ownership marker.** After a successful clear, `Clear` drops a hidden sentinel file
(`.ttsmm-managed`, the `managedMarker` constant) into the directory. On the next run that marker is
proof the tool created the directory, so deleting it is safe. `Clear` refuses unless one of:

- the directory does not exist,
- the directory is empty,
- the directory already contains the marker, or
- the directory has no marker but every file passes the legacy extension allowlist
  (`allowedExtensions`: `.json`, `.gmnotes`, `.luascriptstate`, `.ttslua`, `.xml`) — backward
  compatibility for trees written before the marker existed; they clear once and gain a marker
  going forward.

A non-empty directory with no marker and unrecognized content (a mistargeted `--moddir` pointed at,
say, `$HOME`) is refused with an error explaining how to proceed. This also resolves the older
`.DS_Store`/`.gitkeep` false-positive, since a tool-created directory carries a marker and the
extension check no longer gates it.

**Path guard.** Independent of contents, `pathGuard` resolves the target to an absolute, cleaned
path and refuses the filesystem root, the user's home directory, any ancestor of home, and
suspiciously shallow single-segment paths — targets that are almost always a typo rather than a mod
tree.

## `conversions.go`

Paired helpers for pulling typed values out of a `types.J`:

- `ForceParseIntoX` — returns an error if the key is missing or the wrong type.
- `TryParseIntoX` — the same, discarding the error, for genuinely optional keys.

Both **delete the key on success.** They are consuming reads: the caller extracts the metadata keys
it understands, and whatever remains in the map is passed through untouched. That is how unknown
TTS fields survive a round trip instead of being dropped — so prefer these helpers over reading the
map directly.

Pick `Force` when absence is a real error and `Try` when it isn't; the two exist so the choice is
visible at the call site.

## Testing

`textops_test.go` and `conversions_test.go` cover this package. `TextOps` is built with injected
`readFileToBytes`/`writeBytesToFile` function fields specifically so tests can substitute them
without touching disk — follow that pattern for new I/O types.
