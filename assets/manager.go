package assets

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const assetsBasePath = "assets"

type Resource[T any] struct{ Data T }

type ResourceLoader[T any] struct {
	Load   func(path string) (T, error)
	Unload func(T)
}

type ResourceCache[T any] struct {
	loader ResourceLoader[T]
	items  map[string]*Resource[T]
	mu     sync.RWMutex
}

func NewResourceCache[T any](loader ResourceLoader[T]) *ResourceCache[T] {
	return &ResourceCache[T]{loader: loader, items: make(map[string]*Resource[T])}
}

func (c *ResourceCache[T]) Load(key string) (*Resource[T], error) {
	c.mu.RLock()
	if r, ok := c.items[key]; ok {
		c.mu.RUnlock()
		return r, nil
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	val, err := c.loader.Load(key)
	if err != nil {
		return nil, err
	}

	res := &Resource[T]{Data: val}
	c.items[key] = res
	return res, nil
}

func (c *ResourceCache[T]) Get(key string) (*Resource[T], bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	r, ok := c.items[key]
	return r, ok
}

func (c *ResourceCache[T]) Reload(key string) error {
	c.mu.RLock()
	old, ok := c.items[key]
	c.mu.RUnlock()
	if !ok {
		return fmt.Errorf("asset not loaded: %s", key)
	}

	newVal, err := c.loader.Load(key)
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.loader.Unload(old.Data)
	old.Data = newVal
	return nil
}

func (c *ResourceCache[T]) Unload(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if r, ok := c.items[key]; ok {
		c.loader.Unload(r.Data)
		delete(c.items, key)
	}
}

func (c *ResourceCache[T]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for k, r := range c.items {
		c.loader.Unload(r.Data)
		delete(c.items, k)
	}
}

type AssetManager struct {
	models   *ResourceCache[rl.Model]
	textures *ResourceCache[rl.Texture2D]
	images   *ResourceCache[rl.Image]
	sounds   *ResourceCache[rl.Sound]
	music    *ResourceCache[rl.Music]
	fonts    *ResourceCache[rl.Font]
	shaders  *ResourceCache[rl.Shader]
}

var GlobalManager *AssetManager

func Init() {
	GlobalManager = NewAssetManager()
}

func NewAssetManager() *AssetManager {
	am := &AssetManager{}

	// Here is the set of actual loading functions, those are used by generic ResourceCache.Load() method

	am.models = NewResourceCache(ResourceLoader[rl.Model]{
		Load: func(key string) (rl.Model, error) {
			path := filepath.Join(assetsBasePath, "models", key)
			m := rl.LoadModel(path)
			if m.MeshCount == 0 {
				return rl.Model{}, fmt.Errorf("failed to load model: %s", key)
			}
			return m, nil
		},
		Unload: func(m rl.Model) { rl.UnloadModel(m) },
	})

	am.textures = NewResourceCache(ResourceLoader[rl.Texture2D]{
		Load: func(key string) (rl.Texture2D, error) {
			path := filepath.Join(assetsBasePath, "textures", key)
			tex := rl.LoadTexture(path)
			if tex.ID == 0 {
				return rl.Texture2D{}, fmt.Errorf("failed to load texture: %s", key)
			}
			return tex, nil
		},
		Unload: func(t rl.Texture2D) { rl.UnloadTexture(t) },
	})

	am.images = NewResourceCache(ResourceLoader[rl.Image]{
		Load: func(key string) (rl.Image, error) {
			path := filepath.Join(assetsBasePath, "images", key)
			img := rl.LoadImage(path)
			if img.Data == nil {
				return rl.Image{}, fmt.Errorf("failed to load image: %s", key)
			}
			return *img, nil
		},
		Unload: func(i rl.Image) { rl.UnloadImage(&i) },
	})

	am.sounds = NewResourceCache(ResourceLoader[rl.Sound]{
		Load: func(key string) (rl.Sound, error) {
			path := filepath.Join(assetsBasePath, "audio", key)
			s := rl.LoadSound(path)
			// Note: raylib Sound detection differs; check Stream.Buffer or ID fields as needed
			return s, nil
		},
		Unload: func(s rl.Sound) { rl.UnloadSound(s) },
	})

	am.music = NewResourceCache(ResourceLoader[rl.Music]{
		Load: func(key string) (rl.Music, error) {
			path := filepath.Join(assetsBasePath, "audio", key)
			m := rl.LoadMusicStream(path)
			return m, nil
		},
		Unload: func(m rl.Music) { rl.UnloadMusicStream(m) },
	})

	am.fonts = NewResourceCache(ResourceLoader[rl.Font]{
		Load: func(key string) (rl.Font, error) {
			parts := strings.SplitN(key, ":", 2)
			path := filepath.Join(assetsBasePath, "fonts", parts[0])
			size := int32(16) // default size, TODO: create a const for that
			if len(parts) == 2 {
				var s int
				fmt.Sscanf(parts[1], "%d", &s)
				size = int32(s)
			}
			f := rl.LoadFontEx(path, size, nil, 0)
			if f.Texture.ID == 0 {
				return rl.Font{}, fmt.Errorf("failed to load font: %s", key)
			}
			return f, nil
		},
		Unload: func(f rl.Font) { rl.UnloadFont(f) },
	})

	am.shaders = NewResourceCache(ResourceLoader[rl.Shader]{
		Load: func(key string) (rl.Shader, error) {
			parts := strings.SplitN(key, "|", 2)
			if len(parts) != 2 {
				return rl.Shader{}, fmt.Errorf("invalid shader key: %s", key)
			}
			vs := filepath.Join(assetsBasePath, "shaders", parts[0])
			fs := filepath.Join(assetsBasePath, "shaders", parts[1])
			s := rl.LoadShader(vs, fs)
			if s.ID == 0 {
				return rl.Shader{}, fmt.Errorf("failed to load shader: %s", key)
			}
			return s, nil
		},
		Unload: func(s rl.Shader) { rl.UnloadShader(s) },
	})

	return am
}

// Model wrappers

func (am *AssetManager) LoadModel(filename string) (*Resource[rl.Model], error) {
	return am.models.Load(filename)
}
func (am *AssetManager) GetModel(filename string) (*Resource[rl.Model], bool) {
	return am.models.Get(filename)
}
func (am *AssetManager) ReloadModel(filename string) error { return am.models.Reload(filename) }
func (am *AssetManager) UnloadModel(filename string)       { am.models.Unload(filename) }

// Texture wrappers

func (am *AssetManager) LoadTexture(filename string) (*Resource[rl.Texture2D], error) {
	return am.textures.Load(filename)
}
func (am *AssetManager) GetTexture(filename string) (*Resource[rl.Texture2D], bool) {
	return am.textures.Get(filename)
}
func (am *AssetManager) ReloadTexture(filename string) error { return am.textures.Reload(filename) }
func (am *AssetManager) UnloadTexture(filename string)       { am.textures.Unload(filename) }

// Image wrappers

func (am *AssetManager) LoadImage(filename string) (*Resource[rl.Image], error) {
	return am.images.Load(filename)
}
func (am *AssetManager) GetImage(filename string) (*Resource[rl.Image], bool) {
	return am.images.Get(filename)
}
func (am *AssetManager) ReloadImage(filename string) error { return am.images.Reload(filename) }
func (am *AssetManager) UnloadImage(filename string)       { am.images.Unload(filename) }

// Sound wrappers

func (am *AssetManager) LoadSound(filename string) (*Resource[rl.Sound], error) {
	return am.sounds.Load(filename)
}
func (am *AssetManager) GetSound(filename string) (*Resource[rl.Sound], bool) {
	return am.sounds.Get(filename)
}
func (am *AssetManager) ReloadSound(filename string) error { return am.sounds.Reload(filename) }
func (am *AssetManager) UnloadSound(filename string)       { am.sounds.Unload(filename) }

// Music wrappers

func (am *AssetManager) LoadMusic(filename string) (*Resource[rl.Music], error) {
	return am.music.Load(filename)
}
func (am *AssetManager) GetMusic(filename string) (*Resource[rl.Music], bool) {
	return am.music.Get(filename)
}
func (am *AssetManager) ReloadMusic(filename string) error { return am.music.Reload(filename) }
func (am *AssetManager) UnloadMusic(filename string)       { am.music.Unload(filename) }

// Font wrappers

func (am *AssetManager) LoadFont(filename string, size int) (*Resource[rl.Font], error) {
	key := fmt.Sprintf("%s:%d", filename, size)
	return am.fonts.Load(key)
}
func (am *AssetManager) GetFont(filename string, size int) (*Resource[rl.Font], bool) {
	key := fmt.Sprintf("%s:%d", filename, size)
	return am.fonts.Get(key)
}
func (am *AssetManager) ReloadFont(filename string, size int) error {
	key := fmt.Sprintf("%s:%d", filename, size)
	return am.fonts.Reload(key)
}
func (am *AssetManager) UnloadFont(filename string, size int) {
	key := fmt.Sprintf("%s:%d", filename, size)
	am.fonts.Unload(key)
}

// Shader wrappers

func (am *AssetManager) LoadShader(vsFile, fsFile string) (*Resource[rl.Shader], error) {
	key := vsFile + "|" + fsFile
	return am.shaders.Load(key)
}
func (am *AssetManager) GetShader(vsFile, fsFile string) (*Resource[rl.Shader], bool) {
	return am.shaders.Get(vsFile + "|" + fsFile)
}
func (am *AssetManager) ReloadShader(vsFile, fsFile string) error {
	return am.shaders.Reload(vsFile + "|" + fsFile)
}
func (am *AssetManager) UnloadShader(vsFile, fsFile string) { am.shaders.Unload(vsFile + "|" + fsFile) }

// Clears managers cache
func (am *AssetManager) ClearAll() {
	am.models.Clear()
	am.textures.Clear()
	am.images.Clear()
	am.sounds.Clear()
	am.music.Clear()
	am.fonts.Clear()
	am.shaders.Clear()
}
