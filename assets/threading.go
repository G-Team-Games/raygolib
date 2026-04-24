package assets

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func currentGID() int64 {
	buf := make([]byte, 64)
	n := runtime.Stack(buf, false)
	line := string(buf[:n])
	line = strings.TrimPrefix(line, "goroutine ")
	idx := strings.IndexByte(line, ' ')
	if idx < 0 {
		return 0
	}
	id, err := strconv.ParseInt(line[:idx], 10, 64)
	if err != nil {
		return 0
	}
	return id
}

func (m *Manager) isOwnerThread() bool {
	return currentGID() == m.ownerGID
}

func (m *Manager) runOrQueue(fn func()) error {
	if m.isOwnerThread() {
		fn()
		return nil
	}

	if m.cfg.ThreadPolicy == ThreadPolicyStrict {
		return fmt.Errorf("assets: off-main-thread mutation blocked by strict thread policy")
	}

	select {
	case m.opQueue <- fn:
		return nil
	default:
		return fmt.Errorf("assets: operation queue is full")
	}
}

func (m *Manager) runOrQueueErr(fn func() error) error {
	if m.isOwnerThread() {
		return fn()
	}

	if m.cfg.ThreadPolicy == ThreadPolicyStrict {
		return fmt.Errorf("assets: off-main-thread mutation blocked by strict thread policy")
	}

	errCh := make(chan error, 1)
	select {
	case m.opQueue <- func() { errCh <- fn() }:
		return <-errCh
	default:
		return fmt.Errorf("assets: operation queue is full")
	}
}

func (m *Manager) Tick() {
	for {
		select {
		case fn := <-m.opQueue:
			if fn != nil {
				fn()
			}
		default:
			return
		}
	}
}

func (m *Manager) LoadTexture(filename string) (*Resource[rl.Texture2D], error) {
	if m.isOwnerThread() {
		return m.AssetManager.LoadTexture(filename)
	}

	if m.cfg.ThreadPolicy == ThreadPolicyStrict {
		return nil, fmt.Errorf("assets: off-main-thread mutation blocked by strict thread policy")
	}

	resultCh := make(chan struct {
		res *Resource[rl.Texture2D]
		err error
	}, 1)
	err := m.runOrQueue(func() {
		res, runErr := m.AssetManager.LoadTexture(filename)
		resultCh <- struct {
			res *Resource[rl.Texture2D]
			err error
		}{res: res, err: runErr}
	})
	if err != nil {
		return nil, err
	}

	out := <-resultCh
	return out.res, out.err
}

func (m *Manager) ReloadTexture(filename string) error {
	return m.runOrQueueErr(func() error {
		return m.AssetManager.ReloadTexture(filename)
	})
}

func (m *Manager) LoadModel(filename string) (*Resource[rl.Model], error) {
	if m.isOwnerThread() {
		return m.AssetManager.LoadModel(filename)
	}
	if m.cfg.ThreadPolicy == ThreadPolicyStrict {
		return nil, fmt.Errorf("assets: off-main-thread mutation blocked by strict thread policy")
	}

	resultCh := make(chan struct {
		res *Resource[rl.Model]
		err error
	}, 1)
	err := m.runOrQueue(func() {
		res, runErr := m.AssetManager.LoadModel(filename)
		resultCh <- struct {
			res *Resource[rl.Model]
			err error
		}{res: res, err: runErr}
	})
	if err != nil {
		return nil, err
	}
	out := <-resultCh
	return out.res, out.err
}

func (m *Manager) ReloadModel(filename string) error {
	return m.runOrQueueErr(func() error {
		return m.AssetManager.ReloadModel(filename)
	})
}

func (m *Manager) UnloadModel(filename string) error {
	return m.runOrQueue(func() {
		m.AssetManager.UnloadModel(filename)
	})
}

func (m *Manager) LoadImage(filename string) (*Resource[rl.Image], error) {
	if m.isOwnerThread() {
		return m.AssetManager.LoadImage(filename)
	}
	if m.cfg.ThreadPolicy == ThreadPolicyStrict {
		return nil, fmt.Errorf("assets: off-main-thread mutation blocked by strict thread policy")
	}

	resultCh := make(chan struct {
		res *Resource[rl.Image]
		err error
	}, 1)
	err := m.runOrQueue(func() {
		res, runErr := m.AssetManager.LoadImage(filename)
		resultCh <- struct {
			res *Resource[rl.Image]
			err error
		}{res: res, err: runErr}
	})
	if err != nil {
		return nil, err
	}
	out := <-resultCh
	return out.res, out.err
}

func (m *Manager) ReloadImage(filename string) error {
	return m.runOrQueueErr(func() error {
		return m.AssetManager.ReloadImage(filename)
	})
}

func (m *Manager) UnloadImage(filename string) error {
	return m.runOrQueue(func() {
		m.AssetManager.UnloadImage(filename)
	})
}

func (m *Manager) LoadSound(filename string) (*Resource[rl.Sound], error) {
	if m.isOwnerThread() {
		return m.AssetManager.LoadSound(filename)
	}
	if m.cfg.ThreadPolicy == ThreadPolicyStrict {
		return nil, fmt.Errorf("assets: off-main-thread mutation blocked by strict thread policy")
	}

	resultCh := make(chan struct {
		res *Resource[rl.Sound]
		err error
	}, 1)
	err := m.runOrQueue(func() {
		res, runErr := m.AssetManager.LoadSound(filename)
		resultCh <- struct {
			res *Resource[rl.Sound]
			err error
		}{res: res, err: runErr}
	})
	if err != nil {
		return nil, err
	}
	out := <-resultCh
	return out.res, out.err
}

func (m *Manager) ReloadSound(filename string) error {
	return m.runOrQueueErr(func() error {
		return m.AssetManager.ReloadSound(filename)
	})
}

func (m *Manager) UnloadSound(filename string) error {
	return m.runOrQueue(func() {
		m.AssetManager.UnloadSound(filename)
	})
}

func (m *Manager) LoadMusic(filename string) (*Resource[rl.Music], error) {
	if m.isOwnerThread() {
		return m.AssetManager.LoadMusic(filename)
	}
	if m.cfg.ThreadPolicy == ThreadPolicyStrict {
		return nil, fmt.Errorf("assets: off-main-thread mutation blocked by strict thread policy")
	}

	resultCh := make(chan struct {
		res *Resource[rl.Music]
		err error
	}, 1)
	err := m.runOrQueue(func() {
		res, runErr := m.AssetManager.LoadMusic(filename)
		resultCh <- struct {
			res *Resource[rl.Music]
			err error
		}{res: res, err: runErr}
	})
	if err != nil {
		return nil, err
	}
	out := <-resultCh
	return out.res, out.err
}

func (m *Manager) ReloadMusic(filename string) error {
	return m.runOrQueueErr(func() error {
		return m.AssetManager.ReloadMusic(filename)
	})
}

func (m *Manager) UnloadMusic(filename string) error {
	return m.runOrQueue(func() {
		m.AssetManager.UnloadMusic(filename)
	})
}

func (m *Manager) LoadFont(filename string, size int) (*Resource[rl.Font], error) {
	if m.isOwnerThread() {
		return m.AssetManager.LoadFont(filename, size)
	}
	if m.cfg.ThreadPolicy == ThreadPolicyStrict {
		return nil, fmt.Errorf("assets: off-main-thread mutation blocked by strict thread policy")
	}

	resultCh := make(chan struct {
		res *Resource[rl.Font]
		err error
	}, 1)
	err := m.runOrQueue(func() {
		res, runErr := m.AssetManager.LoadFont(filename, size)
		resultCh <- struct {
			res *Resource[rl.Font]
			err error
		}{res: res, err: runErr}
	})
	if err != nil {
		return nil, err
	}
	out := <-resultCh
	return out.res, out.err
}

func (m *Manager) ReloadFont(filename string, size int) error {
	return m.runOrQueueErr(func() error {
		return m.AssetManager.ReloadFont(filename, size)
	})
}

func (m *Manager) UnloadFont(filename string, size int) error {
	return m.runOrQueue(func() {
		m.AssetManager.UnloadFont(filename, size)
	})
}

func (m *Manager) LoadShader(vsFile, fsFile string) (*Resource[rl.Shader], error) {
	if m.isOwnerThread() {
		return m.AssetManager.LoadShader(vsFile, fsFile)
	}
	if m.cfg.ThreadPolicy == ThreadPolicyStrict {
		return nil, fmt.Errorf("assets: off-main-thread mutation blocked by strict thread policy")
	}

	resultCh := make(chan struct {
		res *Resource[rl.Shader]
		err error
	}, 1)
	err := m.runOrQueue(func() {
		res, runErr := m.AssetManager.LoadShader(vsFile, fsFile)
		resultCh <- struct {
			res *Resource[rl.Shader]
			err error
		}{res: res, err: runErr}
	})
	if err != nil {
		return nil, err
	}
	out := <-resultCh
	return out.res, out.err
}

func (m *Manager) ReloadShader(vsFile, fsFile string) error {
	return m.runOrQueueErr(func() error {
		return m.AssetManager.ReloadShader(vsFile, fsFile)
	})
}

func (m *Manager) UnloadShader(vsFile, fsFile string) error {
	return m.runOrQueue(func() {
		m.AssetManager.UnloadShader(vsFile, fsFile)
	})
}

func (m *Manager) UnloadTexture(filename string) error {
	return m.runOrQueue(func() {
		m.AssetManager.UnloadTexture(filename)
	})
}

func (m *Manager) ClearAll() error {
	return m.runOrQueue(func() {
		m.AssetManager.ClearAll()
	})
}

func (m *Manager) Close() error {
	return m.ClearAll()
}
