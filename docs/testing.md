# Testing guide for _raygolib_

This document groups tests into categories by package and describes how to run them.

## Categories
- Core runtime loop and middleware (`rgl`): 
  - `game_test.go`
  - `debug_test.go`
- Assets manager and asset generator (`rga`):
  - `assets/manager_test.go`
  - `assets/generate_test.go`
  - `assets/asset_kind_test.go`
- Physics 3D collision 3D (`rgcol3d`): 
  - `physics/collision/3d/collision_test.go`
  - `physics/collision/3d/resolver_test.go`
  - `physics/collision/3d/draw_colliders_test.go`

## Detailed catalogs
- Core runtime: `docs/tests/core.md`
- Assets: `docs/tests/assets.md`
- Collision 3D: `docs/tests/collision3d.md`

## Run tests
- Full suite: `make test`
- Coverage (filters `internal/raylib`, `internal/testutils` and `dev/scripts`): `make coverage`

### Run by package (directory)
- Core runtime package tests:
  - `go test -v ./`
- Assets tests:
  - `go test -v ./assets`
- Collision 3D tests:
  - `go test -v ./physics/collision/3d`

## Conventions for new tests
- Tests live in the same package as the code they test, with `_test.go` suffix.
- Use suite-level naming pattern for groupability, for example `TestRunSuite` with `t.Run("loop/...", ...)`.
- Add clear contract wording for invariants (`symmetry`, `order independence`, `zero penetration on touch`).
- Update corresponding `docs/tests/*.md` file when adding/removing test.
