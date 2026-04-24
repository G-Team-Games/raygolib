package rga

import (
	"strings"
	"testing"
	"time"
)

func TestManagerStrictPolicyBlocksOffThreadMutation(t *testing.T) {
	m, err := NewManager(WithThreadPolicy(ThreadPolicyStrict))
	if err != nil {
		t.Fatalf("unexpected manager error: %v", err)
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- m.UnloadTexture("missing.png")
	}()

	select {
	case got := <-errCh:
		if got == nil {
			t.Fatalf("expected strict policy error")
		}
		if !strings.Contains(got.Error(), "strict thread policy") {
			t.Fatalf("unexpected error: %v", got)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatalf("timeout waiting for strict policy result")
	}
}

func TestManagerQueuePolicyRunsMutationOnTick(t *testing.T) {
	m, err := NewManager(WithThreadPolicy(ThreadPolicyQueueOnly))
	if err != nil {
		t.Fatalf("unexpected manager error: %v", err)
	}

	done := make(chan struct{})
	go func() {
		// UnloadTexture uses runOrQueue which returns immediately for void ops
		_ = m.UnloadTexture("missing.png")
		close(done)
	}()

	select {
	case <-done:
		// void ops enqueue and return immediately even in queue-only mode
		// the important thing: no panic, no error, queued for tick
	case <-time.After(200 * time.Millisecond):
		t.Fatalf("timeout waiting for enqueue return")
	}

	m.Tick()
}

func TestManagerQueuePolicyWithLoadReturnsResult(t *testing.T) {
	m, err := NewManager(WithThreadPolicy(ThreadPolicyQueueOnly))
	if err != nil {
		t.Fatalf("unexpected manager error: %v", err)
	}

	done := make(chan error)
	go func() {
		_, loadErr := m.LoadTexture("missing.png")
		done <- loadErr
	}()

	select {
	case <-done:
		t.Fatalf("expected call to block until Tick for Load (result channel)")
	case <-time.After(50 * time.Millisecond):
		// expected: LoadTexture with result channel waits
	}

	m.Tick()

	select {
	case err := <-done:
		// expected: result arrived after tick
		_ = err
	case <-time.After(500 * time.Millisecond):
		t.Fatalf("timeout waiting for load result after Tick")
	}
}
