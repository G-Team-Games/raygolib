package raygolib

import (
	"errors"
	"testing"
)

type TestGame struct {
	updateCalls int
}

func (g *TestGame) Update(dt float32) error {
	g.updateCalls++
	return nil
}

func (g *TestGame) Draw() {}

func TestCallingUpdate(t *testing.T) {
	game := &TestGame{}
	err := update(game, 0.016)

	if err != nil {
		t.Fatal(err)
	}

	if game.updateCalls != 1 {
		t.Fatalf("expected 1 update call, got %d", game.updateCalls)
	}

}

type ErrorGame struct{}

func (g *ErrorGame) Update(dt float32) error {
	return errors.New("update failed")
}

func (g *ErrorGame) Draw() {}

func TestUpdateReturnsError(t *testing.T) {
	game := &ErrorGame{}
	err := update(game, 0.016)

	if err == nil {
		t.Fatal("expected error")
	}
}
