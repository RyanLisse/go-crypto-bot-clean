package middleware

type mockLogger struct {
	infoCalled  bool
	errorCalled bool
	infoArgs    []interface{}
	errorArgs   []interface{}
}

func (m *mockLogger) Info(args ...interface{}) {
	m.infoCalled = true
	m.infoArgs = args
}

func (m *mockLogger) Error(args ...interface{}) {
	m.errorCalled = true
	m.errorArgs = args
}
