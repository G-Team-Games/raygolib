package assets

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestResourceCacheLoadDeduplicatesConcurrent(t *testing.T) {
	var loadCalls int32

	cache := NewResourceCache(ResourceLoader[int]{
		Load: func(key string) (int, error) {
			atomic.AddInt32(&loadCalls, 1)
			time.Sleep(20 * time.Millisecond)
			return 42, nil
		},
		Unload: func(int) {},
	})

	const workers = 16
	results := make([]*Resource[int], workers)
	errs := make([]error, workers)

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			results[idx], errs[idx] = cache.Load("same")
		}(i)
	}
	wg.Wait()

	if got := atomic.LoadInt32(&loadCalls); got != 1 {
		t.Fatalf("expected 1 load call, got %d", got)
	}

	for i := range errs {
		if errs[i] != nil {
			t.Fatalf("unexpected error at %d: %v", i, errs[i])
		}
		if results[i] == nil || results[i].Data != 42 {
			t.Fatalf("unexpected result at %d: %+v", i, results[i])
		}
	}
}

func TestResourceCacheLoadRetryAfterFailure(t *testing.T) {
	var attempts int32

	cache := NewResourceCache(ResourceLoader[string]{
		Load: func(key string) (string, error) {
			curr := atomic.AddInt32(&attempts, 1)
			if curr == 1 {
				return "", errors.New("boom")
			}
			return "ok", nil
		},
		Unload: func(string) {},
	})

	if _, err := cache.Load("retry"); err == nil {
		t.Fatalf("expected first load to fail")
	}

	res, err := cache.Load("retry")
	if err != nil {
		t.Fatalf("expected second load to succeed, got %v", err)
	}
	if res.Data != "ok" {
		t.Fatalf("unexpected value: %q", res.Data)
	}
}

func TestResourceCacheReloadFailureKeepsOldValue(t *testing.T) {
	var loadCalls int32
	var unloadCalls int32

	cache := NewResourceCache(ResourceLoader[int]{
		Load: func(key string) (int, error) {
			curr := atomic.AddInt32(&loadCalls, 1)
			if curr == 1 {
				return 7, nil
			}
			return 0, errors.New("reload failed")
		},
		Unload: func(v int) {
			_ = v
			atomic.AddInt32(&unloadCalls, 1)
		},
	})

	res, err := cache.Load("asset")
	if err != nil {
		t.Fatalf("unexpected load error: %v", err)
	}

	if err := cache.Reload("asset"); err == nil {
		t.Fatalf("expected reload error")
	}

	if res.Data != 7 {
		t.Fatalf("expected old value to stay, got %d", res.Data)
	}

	if got := atomic.LoadInt32(&unloadCalls); got != 0 {
		t.Fatalf("expected no unload on failed reload, got %d", got)
	}
}

func TestResourceCacheReloadSuccessSwapsAndUnloadsOld(t *testing.T) {
	var loadCalls int32
	var unloaded []int
	var mu sync.Mutex

	cache := NewResourceCache(ResourceLoader[int]{
		Load: func(key string) (int, error) {
			curr := atomic.AddInt32(&loadCalls, 1)
			if curr == 1 {
				return 10, nil
			}
			return 20, nil
		},
		Unload: func(v int) {
			mu.Lock()
			defer mu.Unlock()
			unloaded = append(unloaded, v)
		},
	})

	res, err := cache.Load("asset")
	if err != nil {
		t.Fatalf("unexpected load error: %v", err)
	}

	if err := cache.Reload("asset"); err != nil {
		t.Fatalf("unexpected reload error: %v", err)
	}

	if res.Data != 20 {
		t.Fatalf("expected swapped value 20, got %d", res.Data)
	}

	mu.Lock()
	defer mu.Unlock()
	if len(unloaded) != 1 || unloaded[0] != 10 {
		t.Fatalf("unexpected unloaded values: %v", unloaded)
	}
}
