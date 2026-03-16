package raygolib

import (
	"errors"
	"testing"
)

type TestGame struct {
	updateCalls int
	renderCalls int
}

func (g *TestGame) Update(dt float32) error {
	g.updateCalls++
	return nil
}

func (g *TestGame) Render() {
	g.renderCalls++
}

func TestStepCallsUpdateAndRender(t *testing.T) {
	game := &TestGame{}
	err := Step(game, 0.016)

	if err != nil {
		t.Fatal(err)
	}

	if game.updateCalls != 1 {
		t.Fatalf("expected 1 update call, got %d", game.updateCalls)
	}

	if game.renderCalls != 1 {
		t.Fatalf("expected 1 render call, got %d", game.renderCalls)
	}
}

type ErrorGame struct{}

func (g *ErrorGame) Update(dt float32) error {
	return errors.New("update failed")
}

func (g *ErrorGame) Render() {}

func TestStepReturnsError(t *testing.T) {
	game := &ErrorGame{}
	err := Step(game, 0.016)

	if err == nil {
		t.Fatal("expected error")
	}
}
