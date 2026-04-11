package mocks

type InitBackendMock struct {
	InitWindowFunc        func(width, height int32, title string)
	CloseWindowFunc       func()
	SetTargetFPSFunc      func(fps int32)
	WindowShouldCloseFunc func() bool
	GetFrameTimeFunc      func() float32
	BeginDrawingFunc      func()
	EndDrawingFunc        func()

	InitWindowCalls        int
	CloseWindowCalls       int
	SetTargetFPSCalls      int
	WindowShouldCloseCalls int
	GetFrameTimeCalls      int
	BeginDrawingCalls      int
	EndDrawingCalls        int
}

func (m *InitBackendMock) InitWindow(width, height int32, title string) {
	m.InitWindowCalls++
	if m.InitWindowFunc != nil {
		m.InitWindowFunc(width, height, title)
	}
}

func (m *InitBackendMock) CloseWindow() {
	m.CloseWindowCalls++
	if m.CloseWindowFunc != nil {
		m.CloseWindowFunc()
	}
}

func (m *InitBackendMock) SetTargetFPS(fps int32) {
	m.SetTargetFPSCalls++
	if m.SetTargetFPSFunc != nil {
		m.SetTargetFPSFunc(fps)
	}
}

func (m *InitBackendMock) WindowShouldClose() bool {
	m.WindowShouldCloseCalls++
	if m.WindowShouldCloseFunc != nil {
		return m.WindowShouldCloseFunc()
	}
	return true
}

func (m *InitBackendMock) GetFrameTime() float32 {
	m.GetFrameTimeCalls++
	if m.GetFrameTimeFunc != nil {
		return m.GetFrameTimeFunc()
	}
	return 0
}

func (m *InitBackendMock) BeginDrawing() {
	m.BeginDrawingCalls++
	if m.BeginDrawingFunc != nil {
		m.BeginDrawingFunc()
	}
}

func (m *InitBackendMock) EndDrawing() {
	m.EndDrawingCalls++
	if m.EndDrawingFunc != nil {
		m.EndDrawingFunc()
	}
}
