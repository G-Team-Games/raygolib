package rgl

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// ---- Kind detection --------------------------------------------------------

var kindExtensions = map[string]AssetKind{}

var builtInKinds = []AssetKind{
	KindModel,
	KindTexture,
	KindImage,
	KindSound,
	KindMusic,
	KindFont,
	KindShader,
}

func init() {
	kindExtensions = make(map[string]AssetKind)
	for _, kind := range builtInKinds {
		for _, ext := range kind.DefaultExtensions() {
			kindExtensions[strings.ToLower(ext)] = kind
		}
	}
}

func DetectKind(path string) AssetKind {
	ext := strings.ToLower(filepath.Ext(path))
	if kind, ok := kindExtensions[ext]; ok {
		return kind
	}
	return KindUnknown
}

// ---- Resource & loader -----------------------------------------------------

type Resource[T any] struct {
	mu       sync.RWMutex
	Data     T
	path     string
	reloader func() error
	unloader func()
}

// safeRead executes fn while holding a read lock on the resource.
// Use this when reading .Data from a goroutine that isn't the main thread.
func (r *Resource[T]) safeRead(fn func(T)) {
	r.SafeRead(fn)
}

func (r *Resource[T]) SafeRead(fn func(T)) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	fn(r.Data)
}

func (r *Resource[T]) Reload() error {
	r.mu.RLock()
	reload := r.reloader
	r.mu.RUnlock()
	if reload == nil {
		return fmt.Errorf("resource handle is detached")
	}
	rl.TraceLog(rl.LogDebug, "Reloading asset: %s", r.path)
	return reload()
}

func (r *Resource[T]) Unload() {
	r.mu.RLock()
	unload := r.unloader
	r.mu.RUnlock()
	if unload != nil {
		rl.TraceLog(rl.LogDebug, "Unloading asset: %s", r.path)
		unload()
	}
}

func (r *Resource[T]) Path() string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.path
}

type ResourceLoader[T any] struct {
	Load   func(path string) (T, error)
	Unload func(T)
}

// ---- Cache -----------------------------------------------------------------

type cacheEntry[T any] struct {
	resource *Resource[T]
	loading  bool
	cond     *sync.Cond
}

type ResourceCache[T any] struct {
	loader    ResourceLoader[T]
	items     map[string]*cacheEntry[T]
	mu        sync.Mutex
	runOnMain func(fn func())
}

func newResourceCache[T any](loader ResourceLoader[T], runOnMain func(func())) *ResourceCache[T] {
	return &ResourceCache[T]{
		loader:    loader,
		items:     make(map[string]*cacheEntry[T]),
		runOnMain: runOnMain,
	}
}

func (c *ResourceCache[T]) newResource(path string, data T) *Resource[T] {
	res := &Resource[T]{Data: data, path: path}
	res.reloader = func() error {
		return c.reloadByHandle(path, res)
	}
	res.unloader = func() {
		c.Unload(path)
	}
	return res
}

func (c *ResourceCache[T]) zeroResource(res *Resource[T]) T {
	var zero T
	res.mu.Lock()
	old := res.Data
	res.Data = zero
	res.mu.Unlock()
	return old
}

// Load returns a cached resource, loading it (on the main thread) if needed.
// Safe to call from any goroutine; blocks until loading completes.
func (c *ResourceCache[T]) Load(path string) (*Resource[T], error) {
	c.mu.Lock()

	if e, ok := c.items[path]; ok {
		for e.loading {
			e.cond.Wait()
		}
		if e.resource != nil {
			res := e.resource
			c.mu.Unlock()
			return res, nil
		}
	}

	e := &cacheEntry[T]{loading: true}
	e.cond = sync.NewCond(&c.mu)
	c.items[path] = e
	c.mu.Unlock()

	var val T
	var loadErr error
	c.runOnMain(func() {
		val, loadErr = c.loader.Load(path)
	})
	if loadErr != nil {
		c.mu.Lock()
		delete(c.items, path)
		e.loading = false
		e.cond.Broadcast()
		c.mu.Unlock()
		return nil, loadErr
	}

	res := c.newResource(path, val)

	c.mu.Lock()
	e.resource = res
	e.loading = false
	e.cond.Broadcast()
	c.mu.Unlock()
	return res, nil
}

// Get returns an already-loaded resource without triggering a load.
func (c *ResourceCache[T]) Get(path string) (*Resource[T], bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	e, ok := c.items[path]
	if !ok || e.loading || e.resource == nil {
		return nil, false
	}
	return e.resource, true
}

// Reload reloads an already-cached resource in place.
func (c *ResourceCache[T]) Reload(path string) error {
	c.mu.Lock()
	e, ok := c.items[path]
	if !ok {
		c.mu.Unlock()
		return fmt.Errorf("asset not loaded: %s", path)
	}
	for e.loading {
		e.cond.Wait()
	}
	if e.resource == nil {
		c.mu.Unlock()
		return fmt.Errorf("asset not loaded: %s", path)
	}
	e.loading = true
	res := e.resource
	c.mu.Unlock()

	var newVal T
	var loadErr error
	c.runOnMain(func() {
		newVal, loadErr = c.loader.Load(path)
	})
	if loadErr != nil {
		c.mu.Lock()
		e.loading = false
		e.cond.Broadcast()
		c.mu.Unlock()
		return loadErr
	}

	c.runOnMain(func() {
		res.mu.Lock()
		oldVal := res.Data
		res.Data = newVal
		res.mu.Unlock()
		c.loader.Unload(oldVal)
	})

	c.mu.Lock()
	e.loading = false
	e.cond.Broadcast()
	c.mu.Unlock()
	return nil
}

func (c *ResourceCache[T]) reloadByHandle(path string, res *Resource[T]) error {
	c.mu.Lock()
	e, existed := c.items[path]
	if existed {
		for e.loading {
			e.cond.Wait()
		}
	} else {
		e = &cacheEntry[T]{resource: res}
		e.cond = sync.NewCond(&c.mu)
		c.items[path] = e
	}
	e.loading = true
	c.mu.Unlock()

	var newVal T
	var loadErr error
	c.runOnMain(func() {
		newVal, loadErr = c.loader.Load(path)
	})
	if loadErr != nil {
		c.mu.Lock()
		if !existed {
			delete(c.items, path)
		}
		e.loading = false
		e.cond.Broadcast()
		c.mu.Unlock()
		return loadErr
	}

	c.runOnMain(func() {
		var oldVal T
		res.mu.Lock()
		oldVal = res.Data
		res.Data = newVal
		res.path = path
		res.mu.Unlock()
		if existed {
			c.loader.Unload(oldVal)
		}
	})

	c.mu.Lock()
	e.resource = res
	e.loading = false
	e.cond.Broadcast()
	c.mu.Unlock()
	return nil
}

// Unload removes a resource from the cache and frees it on the main thread.
func (c *ResourceCache[T]) Unload(path string) {
	c.mu.Lock()
	e, ok := c.items[path]
	if !ok {
		c.mu.Unlock()
		return
	}
	for e.loading {
		e.cond.Wait()
	}
	delete(c.items, path)
	res := e.resource
	c.mu.Unlock()

	if res != nil {
		c.runOnMain(func() {
			old := c.zeroResource(res)
			c.loader.Unload(old)
		})
	}
}

// Keys returns paths of all fully-loaded entries.
func (c *ResourceCache[T]) Keys() []string {
	c.mu.Lock()
	defer c.mu.Unlock()

	keys := make([]string, 0, len(c.items))
	for k, e := range c.items {
		if e.resource != nil && !e.loading {
			keys = append(keys, k)
		}
	}
	return keys
}

// Clear unloads every entry on the main thread and empties the cache.
func (c *ResourceCache[T]) Clear() {
	c.mu.Lock()
	toUnload := make([]T, 0, len(c.items))
	for k, e := range c.items {
		for e.loading {
			e.cond.Wait()
		}
		if e.resource != nil {
			toUnload = append(toUnload, c.zeroResource(e.resource))
		}
		delete(c.items, k)
	}
	c.mu.Unlock()

	c.runOnMain(func() {
		for _, item := range toUnload {
			c.loader.Unload(item)
		}
	})
}

// ---- Manager ---------------------------------------------------------------

type Manager struct {
	models   *ResourceCache[rl.Model]
	textures *ResourceCache[rl.Texture2D]
	images   *ResourceCache[rl.Image]
	sounds   *ResourceCache[rl.Sound]
	music    *ResourceCache[rl.Music]
	fonts    *ResourceCache[rl.Font]
	shaders  *ResourceCache[rl.Shader]

	opQueue   chan func()
	ownerGID  atomic.Uint64
	closed    atomic.Bool
	closeOnce sync.Once
}

const defaultOpQueueSize = 4096

// NewManager creates a Manager. Must be called from the main/raylib thread.
func NewManager() *Manager {
	m := &Manager{opQueue: make(chan func(), defaultOpQueueSize)}

	runtime.LockOSThread()
	m.ownerGID.Store(currentGoroutineID())

	m.textures = newResourceCache(ResourceLoader[rl.Texture2D]{
		Load: func(path string) (rl.Texture2D, error) {
			tex := rl.LoadTexture(path)
			if tex.ID == 0 {
				return rl.Texture2D{}, fmt.Errorf("failed to load texture: %s", path)
			}
			return tex, nil
		},
		Unload: rl.UnloadTexture,
	}, m.runOnMain)

	m.models = newResourceCache(ResourceLoader[rl.Model]{
		Load: func(path string) (rl.Model, error) {
			model := rl.LoadModel(path)
			if model.MeshCount == 0 {
				return rl.Model{}, fmt.Errorf("failed to load model: %s", path)
			}
			return model, nil
		},
		Unload: rl.UnloadModel,
	}, m.runOnMain)

	m.images = newResourceCache(ResourceLoader[rl.Image]{
		Load: func(path string) (rl.Image, error) {
			img := rl.LoadImage(path)
			if img == nil || img.Data == nil {
				return rl.Image{}, fmt.Errorf("failed to load image: %s", path)
			}
			return *img, nil
		},
		Unload: func(i rl.Image) { rl.UnloadImage(&i) },
	}, m.runOnMain)

	m.sounds = newResourceCache(ResourceLoader[rl.Sound]{
		Load: func(path string) (rl.Sound, error) {
			snd := rl.LoadSound(path)
			if snd.FrameCount == 0 || snd.Stream.Buffer == nil {
				return rl.Sound{}, fmt.Errorf("failed to load sound: %s", path)
			}
			return snd, nil
		},
		Unload: rl.UnloadSound,
	}, m.runOnMain)

	m.music = newResourceCache(ResourceLoader[rl.Music]{
		Load: func(path string) (rl.Music, error) {
			mus := rl.LoadMusicStream(path)
			if mus.FrameCount == 0 || mus.Stream.Buffer == nil {
				return rl.Music{}, fmt.Errorf("failed to load music: %s", path)
			}
			return mus, nil
		},
		Unload: rl.UnloadMusicStream,
	}, m.runOnMain)

	m.fonts = newResourceCache(ResourceLoader[rl.Font]{
		Load: func(path string) (rl.Font, error) {
			parts := strings.SplitN(path, ":", 2)
			fontPath := parts[0]
			size := int32(16)
			if len(parts) == 2 {
				parsed, err := strconv.ParseInt(parts[1], 10, 32)
				if err != nil || parsed <= 0 {
					return rl.Font{}, fmt.Errorf("invalid font size in key: %s", path)
				}
				size = int32(parsed)
			}
			font := rl.LoadFontEx(fontPath, size, nil, 0)
			if font.Texture.ID == 0 {
				return rl.Font{}, fmt.Errorf("failed to load font: %s", path)
			}
			return font, nil
		},
		Unload: rl.UnloadFont,
	}, m.runOnMain)

	m.shaders = newResourceCache(ResourceLoader[rl.Shader]{
		Load: func(path string) (rl.Shader, error) {
			parts := strings.SplitN(path, "|", 2)
			if len(parts) != 2 {
				return rl.Shader{}, fmt.Errorf("invalid shader key (want 'vert|frag'): %s", path)
			}
			s := rl.LoadShader(parts[0], parts[1])
			if s.ID == 0 {
				return rl.Shader{}, fmt.Errorf("failed to load shader: %s", path)
			}
			return s, nil
		},
		Unload: rl.UnloadShader,
	}, m.runOnMain)

	return m
}

// runOnMain executes fn on the owner (main) thread.
// If called from another goroutine, the callback is queued and the caller blocks
// until Tick drains the queue and executes it.
func (m *Manager) runOnMain(fn func()) {
	if fn == nil || m.closed.Load() {
		return
	}

	if m.onOwnerThread() {
		fn()
		return
	}

	done := make(chan struct{})
	queued := true
	func() {
		defer func() {
			if recover() != nil {
				queued = false
			}
		}()
		m.opQueue <- func() {
			fn()
			close(done)
		}
	}()
	if !queued {
		return
	}
	<-done
}

// Tick drains the pending-op queue.
func (m *Manager) Tick() {
	for {
		select {
		case fn, ok := <-m.opQueue:
			if !ok {
				return
			}
			if fn != nil {
				fn()
			}
		default:
			return
		}
	}
}

// ---- Textures --------------------------------------------------------------

func (m *Manager) GetTexture(path string) (*Resource[rl.Texture2D], error) {
	return m.textures.Load(path)
}
func (m *Manager) Texture(path string) *Resource[rl.Texture2D] {
	res, _ := m.textures.Get(path)
	return res
}
func (m *Manager) ReloadTexture(path string) error { return m.textures.Reload(path) }
func (m *Manager) UnloadTexture(path string)       { m.textures.Unload(path) }

// ---- Models ----------------------------------------------------------------

func (m *Manager) GetModel(path string) (*Resource[rl.Model], error) { return m.models.Load(path) }
func (m *Manager) Model(path string) *Resource[rl.Model] {
	res, _ := m.models.Get(path)
	return res
}
func (m *Manager) ReloadModel(path string) error { return m.models.Reload(path) }
func (m *Manager) UnloadModel(path string)       { m.models.Unload(path) }

// ---- Images ----------------------------------------------------------------

func (m *Manager) GetImage(path string) (*Resource[rl.Image], error) { return m.images.Load(path) }
func (m *Manager) Image(path string) *Resource[rl.Image] {
	res, _ := m.images.Get(path)
	return res
}
func (m *Manager) ReloadImage(path string) error { return m.images.Reload(path) }
func (m *Manager) UnloadImage(path string)       { m.images.Unload(path) }

// ---- Sounds ----------------------------------------------------------------

func (m *Manager) GetSound(path string) (*Resource[rl.Sound], error) { return m.sounds.Load(path) }
func (m *Manager) Sound(path string) *Resource[rl.Sound] {
	res, _ := m.sounds.Get(path)
	return res
}
func (m *Manager) ReloadSound(path string) error { return m.sounds.Reload(path) }
func (m *Manager) UnloadSound(path string)       { m.sounds.Unload(path) }

// ---- Music -----------------------------------------------------------------

func (m *Manager) GetMusic(path string) (*Resource[rl.Music], error) { return m.music.Load(path) }
func (m *Manager) Music(path string) *Resource[rl.Music] {
	res, _ := m.music.Get(path)
	return res
}
func (m *Manager) ReloadMusic(path string) error { return m.music.Reload(path) }
func (m *Manager) UnloadMusic(path string)       { m.music.Unload(path) }

// ---- Fonts -----------------------------------------------------------------

func (m *Manager) GetFont(path string) (*Resource[rl.Font], error) { return m.fonts.Load(path) }
func (m *Manager) Font(path string) *Resource[rl.Font] {
	res, _ := m.fonts.Get(path)
	return res
}
func (m *Manager) ReloadFont(path string) error { return m.fonts.Reload(path) }
func (m *Manager) UnloadFont(path string)       { m.fonts.Unload(path) }

// ---- Shaders ---------------------------------------------------------------

func (m *Manager) GetShader(path string) (*Resource[rl.Shader], error) { return m.shaders.Load(path) }
func (m *Manager) Shader(path string) *Resource[rl.Shader] {
	res, _ := m.shaders.Get(path)
	return res
}
func (m *Manager) ReloadShader(path string) error { return m.shaders.Reload(path) }
func (m *Manager) UnloadShader(path string)       { m.shaders.Unload(path) }

// ---- Bulk ops --------------------------------------------------------------

func (m *Manager) Keys(kind AssetKind) []string {
	switch kind {
	case KindModel:
		return m.models.Keys()
	case KindTexture:
		return m.textures.Keys()
	case KindImage:
		return m.images.Keys()
	case KindSound:
		return m.sounds.Keys()
	case KindMusic:
		return m.music.Keys()
	case KindFont:
		return m.fonts.Keys()
	case KindShader:
		return m.shaders.Keys()
	default:
		return nil
	}
}

// ReloadAll reloads all currently loaded assets.
// If onErr is nil, reload errors are ignored.
func (m *Manager) ReloadAll(onErr func(path string, err error)) {
	reloadCache := func(label string, keys []string, reload func(string) error) {
		for _, key := range keys {
			if err := reload(key); err != nil && onErr != nil {
				onErr(key, fmt.Errorf("%s %s: %w", label, key, err))
			}
		}
	}

	for _, c := range m.cacheOrder() {
		reloadCache(c.label, c.keys(), c.reload)
	}
}

func (m *Manager) ClearAll() {
	for _, c := range m.cacheOrder() {
		c.clear()
	}
}

func (m *Manager) Close() error {
	m.closeOnce.Do(func() {
		m.ClearAll()
		m.closed.Store(true)
		close(m.opQueue)
	})
	return nil
}

// ---- Internal helpers ------------------------------------------------------

func (m *Manager) onOwnerThread() bool {
	owner := m.ownerGID.Load()
	if owner == 0 {
		return false
	}
	return currentGoroutineID() == owner
}

func currentGoroutineID() uint64 {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	prefix := "goroutine "
	line := string(buf[:n])
	if !strings.HasPrefix(line, prefix) {
		return 0
	}
	line = line[len(prefix):]
	space := strings.IndexByte(line, ' ')
	if space <= 0 {
		return 0
	}
	id, err := strconv.ParseUint(line[:space], 10, 64)
	if err != nil {
		return 0
	}
	return id
}

type managerCacheOps struct {
	label  string
	keys   func() []string
	reload func(string) error
	clear  func()
}

func (m *Manager) cacheOrder() []managerCacheOps {
	return []managerCacheOps{
		{label: "model", keys: m.models.Keys, reload: m.models.Reload, clear: m.models.Clear},
		{label: "texture", keys: m.textures.Keys, reload: m.textures.Reload, clear: m.textures.Clear},
		{label: "image", keys: m.images.Keys, reload: m.images.Reload, clear: m.images.Clear},
		{label: "sound", keys: m.sounds.Keys, reload: m.sounds.Reload, clear: m.sounds.Clear},
		{label: "music", keys: m.music.Keys, reload: m.music.Reload, clear: m.music.Clear},
		{label: "font", keys: m.fonts.Keys, reload: m.fonts.Reload, clear: m.fonts.Clear},
		{label: "shader", keys: m.shaders.Keys, reload: m.shaders.Reload, clear: m.shaders.Clear},
	}
}
