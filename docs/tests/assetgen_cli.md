# Asset generator CLI tests

Scope: command-line wiring in `cmd/assetgen/main.go` — flag parsing, config mapping, exit codes, stdout/stderr contract, and `parseStringKinds` behavior.

## Kind parsing (`cmd/assetgen/main_test.go`)

- `TestParseStringKinds/empty input`: verifies empty kinds string returns `nil`.
- `TestParseStringKinds/trims and skips empties`: verifies comma/whitespace cleanup and per-kind derived fields (`Dir`, `Type`, `Plural`, `Extensions`).
- `TestParseStringKinds/single letter kind`: verifies one-char kind derives type name correctly (`XName`).

## CLI config mapping (`cmd/assetgen/main_test.go`)

- `TestRun_FlagMappingAndSuccessMessage`: verifies full flag set maps into `rga.Config` and prints success output when generation succeeds.
- `TestRun_KindsParsingOrderPreserved`: verifies `-kinds` order preserved in resulting config.

## Exit code and output contract (`cmd/assetgen/main_test.go`)

- `TestRun_DryRunMessage`: verifies `-dry-run` sets `cfg.DryRun`, prints `Dry run complete`, and avoids normal success message.
- `TestRun_GenerateError`: verifies generator error path returns exit code `1` and prints `Error: ...` to stderr.
- `TestRun_InvalidFlag`: verifies parse failures return exit code `2` and print flag error to stderr.

## Wiring smoke coverage (`cmd/assetgen/main_test.go`)

- `TestRun_DefaultGenerateWhenNil`: verifies nil `generate` fallback uses real `rga.Generate` and produces generated file for minimal temp asset setup.

## Why this category matters

- Protects end-user CLI behavior (flags, outputs, exit codes) independent from internal generator logic.
- Catches regressions when refactoring argument parsing or output stream handling.
