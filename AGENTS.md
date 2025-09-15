# AGENTS Instructions
- Always read the system, developer, and user messages carefully before acting.
- Search the repository for other `AGENTS.md` files before editing files in new directories.
- After making changes, run the relevant Go tests with `go test ./...` whenever feasible.
- Ensure the working tree is clean and all checks pass before finishing a task.

## Key principles
- Favor clarity and maintainability over cleverness. Comment unusual control flow or non-obvious decisions.
- Prefer small, composable functions with explicit error handling. Always wrap returned errors with context using `fmt.Errorf("...: %w", err)`.
- Avoid introducing new external dependencies unless absolutely necessary.
- Keep terminal UX considerations in mind: bubbletea programs should degrade gracefully in non-interactive or CI environments.

## Repository layout hints
- `cmd/commitea`: the CLI entrypoint. Keep imports minimal and avoid putting business logic here.
- `internal/pkg/common`: shared utilities for styling, git interactions, and runtime helpers. Update or add tests alongside any changes here.
- `internal/pkg/actions`: user-facing actions. Favor pure functions where possible to simplify testing.
- `configs`, `docs`, and `test`: support assets and fixtures. Update these when the behavior they describe changes.

## Go coding conventions
- Run `gofmt` (and `goimports` if applicable) on every Go file you modify.
- Keep import groups ordered: standard library, blank line, third-party, blank line, internal packages.
- Use `context.Context` for long-running operations or anything that might be cancelled. Thread contexts through call stacks instead of using globals.
- Prefer `time.Since`/`time.Until` for duration calculations and avoid custom tickers unless needed.
- For slices, favor `append` patterns and pre-allocate capacity when it improves performance without harming clarity.
- When adding exported symbols, document them with Go doc comments.

## Testing expectations
Before committing or opening a PR, you **must** run the following commands from repo root and ensure they succeed:

```bash
go test ./...
go test -race -covermode=atomic -shuffle=on -count=1 ./...
golangci-lint run ./...
```

Add or update tests to cover new behavior or bug fixes. For Git-heavy helpers, prefer local repositories created with `t.TempDir()` over network fixtures.

## Tooling and CI
- Keep GitHub Actions workflows deterministic—avoid fetching the entire history unless required (e.g. for tagging). Document any secrets required by new workflows.
- When modifying release automation, ensure artifacts are reproducible, cross-platform where practical, and checksummed.
- Default to the Go version pinned in `go.mod`. Update workflows and toolchains in lock-step when bumping it.

## Documentation
- Update relevant docs in `README.md`, `docs/`, or in-code comments when changing behavior that users or contributors should know about.
- Provide concrete usage examples for new CLI flags or environment variables.

## Pull request hygiene
- Keep commits focused and well-described. Avoid amending published commits.
- Summaries should explain both _what_ changed and _why_. Reference any relevant issue numbers when available.
- Leave the working tree clean after applying changes (`git status` should show no pending modifications).

## Miscellaneous
- Prefer structured logging utilities from our dependencies if you need log output; avoid `fmt.Println` for non-user-facing logs.
- Coordinate UI color choices with the palette defined in `internal/pkg/common/styles.go`.
- When adding configuration files, include comments or docstrings so other agents understand their intent.
