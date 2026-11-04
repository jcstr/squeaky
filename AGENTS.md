# Squeaky — agent instructions

## What this project is

`squeaky` is a Go CLI for Arch Linux that removes old pacman cache packages,
orphaned packages, old journal logs, stale `~/.cache` files, and old `/tmp`
files. `dry` is read-only. `clean` is the only mutating command.

## Commands

- Build: `make build`
- Test: `make test` (`go test ./... -race`)
- Vet: `go vet ./...`
- Lint: `golangci-lint run ./...`
- Vulnerability scan: `govulncheck ./...`
- Dependency audit: `go list -m -u all`

## Layout

- `cmd/` — Cobra commands. `dry.go` and `clean.go` each build the cleaner list.
- `internal/cleaner/` — one cleaner per file. Every cleaner implements:
  - `Name() string`
  - `Analyze() (*Result, error)` — must never modify the system
  - `Clean() (*Result, error)` — performs the mutation
- `internal/config/` — YAML config and defaults.
- `internal/ui/` — terminal output and summaries.
- `squeaky.example.yaml` — documented user config template.

## Conventions

- `Analyze` must be side-effect-free. Only `Clean` may modify the host.
- Use the package-level `RunCommand` variable for external commands.
  Tests must replace it and restore it with `defer`.
- Add table-driven tests for parsers and fake-command behavior.
- Register a new cleaner in BOTH `cmd/dry.go` and `cmd/clean.go`.
- Put user-configurable values in `internal/config` and document them in
  `squeaky.example.yaml`.
- Prefer the standard library. Add dependencies only with a clear need.
- Do not commit the built `squeaky` binary.

## Safety — non-negotiable

- Never run `squeaky clean`, `paccache -r`, `pacman -Rns`,
  `journalctl --vacuum-*`, or any `sudo` command from an agent session.
- Use unit tests, fake `RunCommand` implementations, or `t.TempDir()`.
- `squeaky dry` is acceptable for read-only inspection, but tests are preferred.
- Do not delete or modify files under `/tmp`, `~/.cache`,
  `/var/cache/pacman`, or the system journal.

## Verification before finishing

- `go test ./... -race`
- `go vet ./...`
- `go build ./...`
- If dependencies changed: `go mod tidy` and `govulncheck ./...`
