---
name: Dependency Audit
description: Check Go modules for outdated, deprecated, or vulnerable dependencies and validate upgrades safely
---

## Workflow

1. Read `go.mod` and note the `go` directive.
2. Run `go list -m -u all` and separate direct from indirect updates.
3. Check whether any dependency requires a newer Go toolchain.
4. Treat `gopkg.in/yaml.v3` as deprecated; use `go.yaml.in/yaml/v3` for YAML v3.
5. Apply changes in a scratch copy or worktree first.
6. Run:
   - `go mod tidy`
   - `go test ./... -race`
   - `go vet ./...`
   - `go build ./...`
   - `govulncheck ./...`
7. Report exact version changes and any toolchain bump.
8. Do not commit unless explicitly asked.
