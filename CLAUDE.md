# Guidance for Claude Code

This repository keeps its architecture notes and working conventions in `AGENTS.md` files — one at
the root and one in each Go package. Those files are the source of truth; this file exists only to
load them into Claude Code's memory.

@AGENTS.md

Each package directory (`bundler/`, `file/`, `handler/`, `mod/`, `objects/`, `tests/`, `types/`)
has its own `AGENTS.md`, pulled in automatically by a small `CLAUDE.md` bridge in that directory
whenever you work on files there. **Read the relevant package's `AGENTS.md` before changing its
code**, and when you change behavior, update that `AGENTS.md` rather than this file.
