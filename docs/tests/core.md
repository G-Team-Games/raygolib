# Core runtime tests

Scope: game initializer/config merge, middleware ordering, game loop orchestration, debug API/middleware behavior.

## `game_test.go`

- `TestInitGameSuite`: verifies `InitGame` stores game + default config, `InitGameWithConfig` merges custom fields while preserving defaults for zero values.
- `TestMergeInitGameConfigsSuite`: verifies non-zero override semantics in `mergeInitGameConfigs`.
- `TestWithMiddlewareSuite`: verifies middleware wrapping order contract (last added is outermost) for init/update/draw/close chain.
- `TestRunSuite`: verifies `Run` integration with backend lifecycle, update/draw call counts, dt forwarding, and update-error early-exit with guaranteed `CloseWindow`.

### Subtests in `TestInitGameSuite`

- `uses default config`: checks initializer keeps passed game and applies canonical defaults.
- `merges custom config`: checks partial custom config overrides only non-zero fields.

### Subtests in `TestMergeInitGameConfigsSuite`

- `overrides non-zero values`: checks full-field non-zero override.

### Subtests in `TestWithMiddlewareSuite`

- `applies in order`: checks call trace order through stacked middleware and wrapped game.

### Subtests in `TestRunSuite`

- `calls game loop and backend`: checks backend init/FPS/drawing/close calls and game update/draw dt values.
- `returns update error and closes window`: checks error propagation, draw short-circuit, and cleanup on failure.

## `debug_test.go`

- `TestDebugAPISuite`: verifies queueing behavior when debug disabled/enabled, flush rendering+clear, and toggle state transitions.
- `TestDebugMiddlewareSuite`: verifies middleware injects global `DebugAPI`, handles toggle key, and draws/skips FPS based on config/state.

### Subtests in `TestDebugAPISuite`

- `Rect disabled does not queue`: no draw operations queued when API disabled.
- `Rect enabled queues`: draw operation queued when API enabled.
- `Flush clears queue and draws`: queued draw commands render to backend and queue clears.
- `Toggle flips enabled state`: toggle switches enabled flag both directions.

### Subtests in `TestDebugMiddlewareSuite`

- `injects and configures`: validates `DebugAware` injection, singleton replacement, wrapper config mapping.
- `Update toggles on key press`: validates key-driven enable toggle on update.
- `Draw calls FPS when enabled`: validates FPS drawing path when enabled and `ShowFPS` true.
- `Draw skips FPS when disabled`: validates no FPS draw and draw queue clear when disabled.
- `Draw skips FPS when hidden`: validates no FPS draw when `ShowFPS` false.
