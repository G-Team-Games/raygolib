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
- Asset generator CLI (`cmd/assetgen`):
  - `cmd/assetgen/main_test.go`
- Physics 3D collision 3D (`rgcol3d`): 
  - `physics/collision/3d/collision_test.go`
  - `physics/collision/3d/resolver_test.go`
  - `physics/collision/3d/draw_colliders_test.go`

## Detailed test docs
- Core runtime: `docs/tests/core.md`
- Assets: `docs/tests/assets.md`
- Asset generator CLI: `docs/tests/assetgen_cli.md`
- Collision 3D: `docs/tests/collision3d.md`

## Run tests
- Full suite: `make test`
- Coverage (filters `internal/raylib`, `internal/testutils` and `dev/scripts`): `make coverage`

### Run by package (directory)
- Core runtime package tests:
  - `go test -v ./`
- Assets tests:
  - `go test -v ./assets`
- Asset generator CLI tests:
  - `go test -v ./cmd/assetgen`
- Collision 3D tests:
  - `go test -v ./physics/collision/3d`

## Conventions for new tests
- Tests live in the same package as the code they test, with `_test.go` suffix
- Update corresponding `docs/tests/*.md` file when adding/removing test
