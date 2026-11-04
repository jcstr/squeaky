---
name: Cleaner Development
description: Add or modify a Squeaky cleaner safely with fake-command tests
---

## Contract

Implement `cleaner.Cleaner`:

- `Name() string`
- `Analyze() (*Result, error)` — read-only
- `Clean() (*Result, error)` — mutating

## Rules

- Use the package-level `RunCommand` variable for external commands.
- In tests, replace `RunCommand` and restore it with `defer`.
- Extract output parsing into pure functions and table-test them.
- Register the cleaner in both `cmd/dry.go` and `cmd/clean.go`.
- Add config fields to `internal/config` and document them in
  `squeaky.example.yaml`.
- Never run destructive host commands in tests; use `t.TempDir()`.
