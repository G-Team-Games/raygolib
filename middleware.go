package raygolib

type Middleware func(Game) Game

func DebugMiddleware() Middleware {
	return func(next Game) Game {
		debug := &DebugAPI{
			enabled: true,
		}

		if g, ok := next.(DebugAware); ok {
			g.SetDebug(debug)
		}

		return &debugWrapper{
			next:  next,
			debug: debug,
		}
	}
}
