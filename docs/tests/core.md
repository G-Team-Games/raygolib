# Core runtime tests

Scope: game initializer/config merge, middleware ordering, game loop orchestration, debug API/middleware behavior.

### Config and init (`game_test.go`)

- `TestInitGameSuite`: verifies `InitGame` stores game + default config, and `InitGameWithConfig` merges custom fields while keeping zero-value defaults.
- `TestInitGameSuite/uses default config`: checks canonical defaults and stored game pointer.
- `TestInitGameSuite/merges custom config`: checks partial override semantics.
- `TestMergeInitGameConfigsSuite`: verifies config merge helper behavior.
- `TestMergeInitGameConfigsSuite/overrides non-zero values`: checks non-zero overrides across all fields.

### Middleware chain (`game_test.go`)

- `TestWithMiddlewareSuite`: verifies wrapping order contract (last added is outermost).
- `TestWithMiddlewareSuite/applies in order`: checks exact init/update/draw/close call sequence.

### Run loop and error path (`game_test.go`)

- `TestRunSuite`: verifies backend lifecycle integration and frame loop behavior.
- `TestRunSuite/calls game loop and backend`: checks window init args, FPS setup, begin/end drawing counts, dt propagation.
- `TestRunSuite/returns update error and closes window`: checks update error propagation, draw short-circuit, cleanup on failure.

### Debug API (`debug_test.go`)

- `TestDebugAPISuite`: verifies draw queueing, flushing, and enabled-state toggling.
- `TestDebugAPISuite/Rect disabled does not queue`: disabled state must not enqueue debug draw ops.
- `TestDebugAPISuite/Rect enabled queues`: enabled state must enqueue draw ops.
- `TestDebugAPISuite/Flush clears queue and draws`: queued ops are rendered then queue is emptied.
- `TestDebugAPISuite/Toggle flips enabled state`: toggle changes state both directions.

### Debug middleware (`debug_test.go`)

- `TestDebugMiddlewareSuite`: verifies injection, key-toggle behavior, and FPS drawing gates.
- `TestDebugMiddlewareSuite/injects and configures`: validates `DebugAware` injection, singleton assignment, wrapper config copy.
- `TestDebugMiddlewareSuite/Update toggles on key press`: validates configured key toggles debug enabled state.
- `TestDebugMiddlewareSuite/Draw calls FPS when enabled`: validates FPS draw path when enabled and `ShowFPS=true`.
- `TestDebugMiddlewareSuite/Draw skips FPS when disabled`: validates no FPS draw and queue clear when disabled.
- `TestDebugMiddlewareSuite/Draw skips FPS when hidden`: validates no FPS draw when `ShowFPS=false`.
