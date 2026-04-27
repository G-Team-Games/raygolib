package rga

import (
	"errors"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func buildTestManager() *Manager {
	runOnMain := func(fn func()) { fn() }
	m := &Manager{}
	m.ownerGID.Store(currentGoroutineID())
	m.opQueue = make(chan func(), 8)

	m.models = newResourceCache(ResourceLoader[rl.Model]{
		Load: func(path string) (rl.Model, error) {
			return rl.Model{MeshCount: 1}, nil
		},
		Unload: func(rl.Model) {},
	}, runOnMain)

	m.textures = newResourceCache(ResourceLoader[rl.Texture2D]{
		Load: func(path string) (rl.Texture2D, error) {
			return rl.Texture2D{ID: 1}, nil
		},
		Unload: func(rl.Texture2D) {},
	}, runOnMain)

	m.images = newResourceCache(ResourceLoader[rl.Image]{
		Load: func(path string) (rl.Image, error) {
			return rl.Image{}, nil
		},
		Unload: func(rl.Image) {},
	}, runOnMain)

	m.sounds = newResourceCache(ResourceLoader[rl.Sound]{
		Load: func(path string) (rl.Sound, error) {
			return rl.Sound{FrameCount: 1}, nil
		},
		Unload: func(rl.Sound) {},
	}, runOnMain)

	m.music = newResourceCache(ResourceLoader[rl.Music]{
		Load: func(path string) (rl.Music, error) {
			return rl.Music{FrameCount: 1}, nil
		},
		Unload: func(rl.Music) {},
	}, runOnMain)

	m.fonts = newResourceCache(ResourceLoader[rl.Font]{
		Load: func(path string) (rl.Font, error) {
			return rl.Font{Texture: rl.Texture2D{ID: 1}}, nil
		},
		Unload: func(rl.Font) {},
	}, runOnMain)

	m.shaders = newResourceCache(ResourceLoader[rl.Shader]{
		Load: func(path string) (rl.Shader, error) {
			return rl.Shader{ID: 1}, nil
		},
		Unload: func(rl.Shader) {},
	}, runOnMain)

	return m
}

type fakeResource struct {
	Value int
}

func TestResourceCacheLoadConcurrentSingleLoader(t *testing.T) {
	var loadCalls atomic.Int32
	release := make(chan struct{})

	cache := newResourceCache(ResourceLoader[fakeResource]{
		Load: func(path string) (fakeResource, error) {
			loadCalls.Add(1)
			<-release
			return fakeResource{Value: 7}, nil
		},
		Unload: func(fakeResource) {},
	}, func(fn func()) { fn() })

	const goroutines = 12
	results := make([]*Resource[fakeResource], goroutines)
	var wg sync.WaitGroup
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			res, err := cache.Load("same-key")
			if err != nil {
				t.Errorf("load failed: %v", err)
				return
			}
			results[idx] = res
		}(i)
	}

	deadline := time.Now().Add(2 * time.Second)
	for loadCalls.Load() == 0 {
		if time.Now().After(deadline) {
			t.Fatal("loader was never called")
		}
		time.Sleep(1 * time.Millisecond)
	}
	if got := loadCalls.Load(); got != 1 {
		t.Fatalf("expected one loader call while waiting, got %d", got)
	}

	close(release)
	wg.Wait()

	first := results[0]
	if first == nil {
		t.Fatal("expected first result")
	}
	for i := 1; i < goroutines; i++ {
		if results[i] != first {
			t.Fatalf("expected all goroutines to share one cached pointer, mismatch at %d", i)
		}
	}
	if first.Data.Value != 7 {
		t.Fatalf("unexpected resource value: %d", first.Data.Value)
	}
}

func TestResourceCacheReloadErrorKeepsOldData(t *testing.T) {
	var calls atomic.Int32
	cache := newResourceCache(ResourceLoader[fakeResource]{
		Load: func(path string) (fakeResource, error) {
			if calls.Add(1) == 1 {
				return fakeResource{Value: 11}, nil
			}
			return fakeResource{}, errors.New("reload failed")
		},
		Unload: func(fakeResource) {},
	}, func(fn func()) { fn() })

	res, err := cache.Load("a")
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	if err := cache.Reload("a"); err == nil {
		t.Fatal("expected reload error")
	}
	if res.Data.Value != 11 {
		t.Fatalf("expected old data to stay after failed reload, got %d", res.Data.Value)
	}
}

func TestResourceHandleUnloadAndReloadReattach(t *testing.T) {
	var value atomic.Int32
	value.Store(1)

	cache := newResourceCache(ResourceLoader[fakeResource]{
		Load: func(path string) (fakeResource, error) {
			return fakeResource{Value: int(value.Load())}, nil
		},
		Unload: func(fakeResource) {},
	}, func(fn func()) { fn() })

	res, err := cache.Load("h")
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	res.Unload()
	if res.Data.Value != 0 {
		t.Fatalf("expected resource to be zeroed after unload, got %d", res.Data.Value)
	}
	if _, ok := cache.Get("h"); ok {
		t.Fatal("expected cache entry to be removed after unload")
	}

	value.Store(22)
	if err := res.Reload(); err != nil {
		t.Fatalf("reload via handle failed: %v", err)
	}

	got, ok := cache.Get("h")
	if !ok {
		t.Fatal("expected cache entry to be reattached")
	}
	if got != res {
		t.Fatal("expected reattached cache entry to use the same handle pointer")
	}
	if res.Data.Value != 22 {
		t.Fatalf("expected reloaded value 22, got %d", res.Data.Value)
	}
	if res.Path() != "h" {
		t.Fatalf("expected path h, got %q", res.Path())
	}

	var seen int
	res.safeRead(func(v fakeResource) {
		seen = v.Value
	})
	if seen != 22 {
		t.Fatalf("safeRead saw %d, want 22", seen)
	}
}

func TestResourceSafeRead(t *testing.T) {
	res := &Resource[fakeResource]{Data: fakeResource{Value: 9}}

	var seen int
	res.SafeRead(func(v fakeResource) {
		seen = v.Value
	})

	if seen != 9 {
		t.Fatalf("SafeRead saw %d, want 9", seen)
	}
}

func TestResourceCacheClearUnloadsAllAndZeros(t *testing.T) {
	var unloadCalls atomic.Int32
	cache := newResourceCache(ResourceLoader[fakeResource]{
		Load: func(path string) (fakeResource, error) {
			if path == "one" {
				return fakeResource{Value: 1}, nil
			}
			return fakeResource{Value: 2}, nil
		},
		Unload: func(fakeResource) {
			unloadCalls.Add(1)
		},
	}, func(fn func()) { fn() })

	one, _ := cache.Load("one")
	two, _ := cache.Load("two")

	cache.Clear()

	if got := unloadCalls.Load(); got != 2 {
		t.Fatalf("expected two unload calls, got %d", got)
	}
	if len(cache.Keys()) != 0 {
		t.Fatal("expected cache keys to be empty after clear")
	}
	if one.Data.Value != 0 || two.Data.Value != 0 {
		t.Fatal("expected loaded handles to be zeroed after clear")
	}
}

func TestResourceHandleReloadFailureWhenDetachedKeepsCacheEmpty(t *testing.T) {
	cache := newResourceCache(ResourceLoader[fakeResource]{
		Load: func(path string) (fakeResource, error) {
			return fakeResource{}, errors.New("cannot load")
		},
		Unload: func(fakeResource) {},
	}, func(fn func()) { fn() })

	res := cache.newResource("k", fakeResource{})
	if err := res.Reload(); err == nil {
		t.Fatal("expected reload error")
	}
	if _, ok := cache.Get("k"); ok {
		t.Fatal("expected failed detached reload to leave cache empty")
	}
}

func TestDetectKind(t *testing.T) {
	tests := []struct {
		path string
		want AssetKind
	}{
		{"a.png", KindTexture},
		{"a.jpg", KindTexture},
		{"a.hdr", KindImage},
		{"a.wav", KindSound},
		{"a.mp3", KindMusic},
		{"a.ttf", KindFont},
		{"a.glsl", KindShader},
		{"a.vs", KindShader},
		{"a.fs", KindShader},
		{"a.obj", KindModel},
		{"A.GLB", KindModel},
		{"a.unknown", KindUnknown},
	}

	for _, tc := range tests {
		if got := DetectKind(tc.path); got != tc.want {
			t.Fatalf("DetectKind(%q) = %q, want %q", tc.path, got, tc.want)
		}
	}
}

func TestManagerDetectKindUsesBuiltInExtensions(t *testing.T) {
	for _, tc := range []struct {
		path string
		want AssetKind
	}{
		{"scene.GLB", KindModel},
		{"sprite.PNG", KindTexture},
		{"image.HDR", KindImage},
		{"sound.MP3", KindMusic},
		{"fx.VS", KindShader},
		{"fx.FS", KindShader},
	} {
		if got := DetectKind(tc.path); got != tc.want {
			t.Fatalf("DetectKind(%q) = %q, want %q", tc.path, got, tc.want)
		}
	}
}

func TestFontLoaderRejectsInvalidSizeKey(t *testing.T) {
	resultCh := make(chan error, 1)
	go func() {
		m := NewManager()
		defer func() { resultCh <- m.Close() }()

		if _, err := m.fonts.loader.Load("font.ttf:not-a-number"); err == nil {
			resultCh <- errors.New("expected invalid font size error")
			return
		}
		if _, err := m.fonts.loader.Load("font.ttf:0"); err == nil {
			resultCh <- errors.New("expected non-positive font size error")
			return
		}
		resultCh <- nil
	}()

	if err := <-resultCh; err != nil {
		t.Fatal(err)
	}
}

func TestManagerReloadAllCallsErrorCallback(t *testing.T) {
	runOnMain := func(fn func()) { fn() }
	m := &Manager{}
	m.ownerGID.Store(currentGoroutineID())
	m.opQueue = make(chan func(), 8)

	var textureCalls atomic.Int32
	m.textures = newResourceCache(ResourceLoader[rl.Texture2D]{
		Load: func(path string) (rl.Texture2D, error) {
			call := textureCalls.Add(1)
			if path == "bad.png" && call > 2 {
				return rl.Texture2D{}, errors.New("boom")
			}
			return rl.Texture2D{ID: 1}, nil
		},
		Unload: func(rl.Texture2D) {},
	}, runOnMain)

	// Remaining caches are needed by ReloadAll and can be no-op successful caches.
	emptyModel := newResourceCache(ResourceLoader[rl.Model]{Load: func(string) (rl.Model, error) { return rl.Model{MeshCount: 1}, nil }, Unload: func(rl.Model) {}}, runOnMain)
	emptyImage := newResourceCache(ResourceLoader[rl.Image]{Load: func(string) (rl.Image, error) { return rl.Image{}, nil }, Unload: func(rl.Image) {}}, runOnMain)
	emptySound := newResourceCache(ResourceLoader[rl.Sound]{Load: func(string) (rl.Sound, error) { return rl.Sound{FrameCount: 1}, nil }, Unload: func(rl.Sound) {}}, runOnMain)
	emptyMusic := newResourceCache(ResourceLoader[rl.Music]{Load: func(string) (rl.Music, error) { return rl.Music{FrameCount: 1}, nil }, Unload: func(rl.Music) {}}, runOnMain)
	emptyFont := newResourceCache(ResourceLoader[rl.Font]{Load: func(string) (rl.Font, error) { return rl.Font{Texture: rl.Texture2D{ID: 1}}, nil }, Unload: func(rl.Font) {}}, runOnMain)
	emptyShader := newResourceCache(ResourceLoader[rl.Shader]{Load: func(string) (rl.Shader, error) { return rl.Shader{ID: 1}, nil }, Unload: func(rl.Shader) {}}, runOnMain)

	m.models = emptyModel
	m.images = emptyImage
	m.sounds = emptySound
	m.music = emptyMusic
	m.fonts = emptyFont
	m.shaders = emptyShader

	if _, err := m.GetTexture("ok.png"); err != nil {
		t.Fatalf("load ok texture failed: %v", err)
	}
	if _, err := m.GetTexture("bad.png"); err != nil {
		t.Fatalf("load bad texture failed: %v", err)
	}

	var paths []string
	var errs []error
	m.ReloadAll(func(path string, err error) {
		paths = append(paths, path)
		errs = append(errs, err)
	})

	if len(paths) != 1 {
		t.Fatalf("expected one callback error, got %d", len(paths))
	}
	if paths[0] != "bad.png" {
		t.Fatalf("unexpected callback path: %q", paths[0])
	}
	if errs[0] == nil || !strings.Contains(errs[0].Error(), "texture bad.png") {
		t.Fatalf("unexpected callback error: %v", errs[0])
	}

	// Ensure nil callback path does not panic.
	m.ReloadAll(nil)
}

func TestManagerWrapperMethodsAndBulkOps(t *testing.T) {
	m := buildTestManager()

	if _, err := m.GetTexture("a.png"); err != nil {
		t.Fatal(err)
	}
	if _, err := m.GetModel("a.obj"); err != nil {
		t.Fatal(err)
	}
	if _, err := m.GetImage("a.hdr"); err != nil {
		t.Fatal(err)
	}
	if _, err := m.GetSound("a.wav"); err != nil {
		t.Fatal(err)
	}
	if _, err := m.GetMusic("a.mp3"); err != nil {
		t.Fatal(err)
	}
	if _, err := m.GetFont("font.ttf:16"); err != nil {
		t.Fatal(err)
	}
	if _, err := m.GetShader("vert.glsl|frag.glsl"); err != nil {
		t.Fatal(err)
	}

	if m.Texture("a.png") == nil || m.Model("a.obj") == nil || m.Image("a.hdr") == nil {
		t.Fatal("expected texture/model/image accessors to return loaded resources")
	}
	if m.Sound("a.wav") == nil || m.Music("a.mp3") == nil || m.Font("font.ttf:16") == nil || m.Shader("vert.glsl|frag.glsl") == nil {
		t.Fatal("expected sound/music/font/shader accessors to return loaded resources")
	}

	if err := m.ReloadTexture("a.png"); err != nil {
		t.Fatal(err)
	}
	if err := m.ReloadModel("a.obj"); err != nil {
		t.Fatal(err)
	}
	if err := m.ReloadImage("a.hdr"); err != nil {
		t.Fatal(err)
	}
	if err := m.ReloadSound("a.wav"); err != nil {
		t.Fatal(err)
	}
	if err := m.ReloadMusic("a.mp3"); err != nil {
		t.Fatal(err)
	}
	if err := m.ReloadFont("font.ttf:16"); err != nil {
		t.Fatal(err)
	}
	if err := m.ReloadShader("vert.glsl|frag.glsl"); err != nil {
		t.Fatal(err)
	}

	if len(m.Keys(KindTexture)) == 0 || len(m.Keys(KindModel)) == 0 || len(m.Keys(KindImage)) == 0 {
		t.Fatal("expected keys for loaded kinds")
	}
	if got := m.Keys(AssetKind("unknown")); got != nil {
		t.Fatal("expected nil keys for unknown kind")
	}

	m.UnloadTexture("a.png")
	m.UnloadModel("a.obj")
	m.UnloadImage("a.hdr")
	m.UnloadSound("a.wav")
	m.UnloadMusic("a.mp3")
	m.UnloadFont("font.ttf:16")
	m.UnloadShader("vert.glsl|frag.glsl")

	if len(m.Keys(KindTexture)) != 0 || len(m.Keys(KindModel)) != 0 || len(m.Keys(KindImage)) != 0 {
		t.Fatal("expected unloaded keys to be empty")
	}

	if _, err := m.GetTexture("b.png"); err != nil {
		t.Fatal(err)
	}
	m.ClearAll()
	if len(m.Keys(KindTexture)) != 0 {
		t.Fatal("expected clear all to remove all keys")
	}

	if err := m.Close(); err != nil {
		t.Fatalf("close failed: %v", err)
	}
}

func TestManagerRunOnMainQueuePath(t *testing.T) {
	m := &Manager{opQueue: make(chan func(), 8)}
	m.ownerGID.Store(^uint64(0))

	done := make(chan struct{})
	go func() {
		m.runOnMain(func() {
			close(done)
		})
	}()

	deadline := time.Now().Add(2 * time.Second)
	for {
		select {
		case <-done:
			goto checked
		default:
		}
		if time.Now().After(deadline) {
			t.Fatal("queued runOnMain callback was not executed by Tick")
		}
		m.Tick()
		time.Sleep(1 * time.Millisecond)
	}

checked:

	if m.onOwnerThread() {
		t.Fatal("expected non-owner thread flag")
	}
}

func TestManagerCloseIsIdempotentAndStopsQueueing(t *testing.T) {
	m := buildTestManager()

	if err := m.Close(); err != nil {
		t.Fatalf("first close failed: %v", err)
	}
	if err := m.Close(); err != nil {
		t.Fatalf("second close failed: %v", err)
	}

	called := false
	m.runOnMain(func() { called = true })
	if called {
		t.Fatal("expected runOnMain to be a no-op after close")
	}
}

func TestManagerKeysSwitchAllKinds(t *testing.T) {
	m := buildTestManager()

	// populate one entry per cache
	_, _ = m.GetTexture("a.png")
	_, _ = m.GetModel("a.obj")
	_, _ = m.GetImage("a.hdr")
	_, _ = m.GetSound("a.wav")
	_, _ = m.GetMusic("a.mp3")
	_, _ = m.GetFont("a.ttf:16")
	_, _ = m.GetShader("v.glsl|f.glsl")

	if len(m.Keys(KindTexture)) != 1 {
		t.Fatal("expected texture keys")
	}
	if len(m.Keys(KindModel)) != 1 {
		t.Fatal("expected model keys")
	}
	if len(m.Keys(KindImage)) != 1 {
		t.Fatal("expected image keys")
	}
	if len(m.Keys(KindSound)) != 1 {
		t.Fatal("expected sound keys")
	}
	if len(m.Keys(KindMusic)) != 1 {
		t.Fatal("expected music keys")
	}
	if len(m.Keys(KindFont)) != 1 {
		t.Fatal("expected font keys")
	}
	if len(m.Keys(KindShader)) != 1 {
		t.Fatal("expected shader keys")
	}
	if got := m.Keys(AssetKind("other")); got != nil {
		t.Fatal("expected nil for unknown kind")
	}
}

func TestNewManagerLoaderClosuresErrorPaths(t *testing.T) {
	// NewManager locks the current goroutine OS thread, so construct it in a helper goroutine.
	mgrCh := make(chan *Manager, 1)
	go func() {
		mgrCh <- NewManager()
	}()
	m := <-mgrCh
	defer func() {
		_ = m.Close()
	}()

	if m == nil || m.textures == nil || m.models == nil || m.images == nil || m.sounds == nil || m.music == nil || m.fonts == nil || m.shaders == nil {
		t.Fatal("expected NewManager to initialize all caches")
	}

	if _, err := m.textures.loader.Load("/definitely/missing_texture.png"); err == nil {
		t.Fatal("expected texture loader error")
	}
	if _, err := m.models.loader.Load("/definitely/missing_model.obj"); err == nil {
		t.Fatal("expected model loader error")
	}
	if _, err := m.images.loader.Load("/definitely/missing_image.png"); err == nil {
		t.Fatal("expected image loader error")
	}

	// Sound and music loaders need audio initialized in many environments.
	rl.InitAudioDevice()
	defer rl.CloseAudioDevice()

	if _, err := m.sounds.loader.Load("/definitely/missing_sound.wav"); err == nil {
		t.Fatal("expected sound loader error")
	}
	if _, err := m.music.loader.Load("/definitely/missing_music.mp3"); err == nil {
		t.Fatal("expected music loader error")
	}
	if _, err := m.fonts.loader.Load("/definitely/missing_font.ttf:16"); err == nil {
		t.Fatal("expected font loader error")
	}

	if _, err := m.shaders.loader.Load("invalid-shader-key"); err == nil {
		t.Fatal("expected shader key format error")
	}
	if _, err := m.shaders.loader.Load("/definitely/missing.vert|/definitely/missing.frag"); err == nil {
		t.Fatal("expected shader loader error")
	}

	called := false
	m.ownerGID.Store(currentGoroutineID())
	m.runOnMain(func() { called = true })
	if !called {
		t.Fatal("expected owner-thread runOnMain to execute immediately")
	}
}
