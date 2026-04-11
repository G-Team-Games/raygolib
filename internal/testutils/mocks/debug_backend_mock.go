package mocks

import "image/color"

type DebugBackendMock struct {
	IsKeyPressedFunc  func(key int32) bool
	DrawFPSFunc       func(x, y int32)
	DrawRectangleFunc func(x, y, w, h int32, color color.RGBA)

	IsKeyPressedCalls  int
	DrawFPSCalls       int
	DrawRectangleCalls int
}

func (m *DebugBackendMock) IsKeyPressed(key int32) bool {
	m.IsKeyPressedCalls++
	if m.IsKeyPressedFunc != nil {
		return m.IsKeyPressedFunc(key)
	}
	return false
}

func (m *DebugBackendMock) DrawFPS(x, y int32) {
	m.DrawFPSCalls++
	if m.DrawFPSFunc != nil {
		m.DrawFPSFunc(x, y)
	}
}

func (m *DebugBackendMock) DrawRectangle(x, y, w, h int32, c color.RGBA) {
	m.DrawRectangleCalls++
	if m.DrawRectangleFunc != nil {
		m.DrawRectangleFunc(x, y, w, h, c)
	}
}
