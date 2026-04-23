# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

`portdor` — a local dev service registry and manager. Single static Go binary with a CLI and HTMX Web UI at `localhost:4242`.

## Git Workflow

Branch protection is enabled on both `dev` and `main`. All changes must follow this flow:

```
feature branch → dev (PR) → main (PR)
```

- Always branch off `dev`: `git checkout dev && git pull && git checkout -b feat/<name>`
- Never commit directly to `dev` or `main`
- When merging PRs, always use **"Create a merge commit"** — squash and rebase break the shared history between `dev` and `main`
