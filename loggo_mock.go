package loggo

type MockLogger struct {
	// You can add more fields if needed
	InfoMessages  []string
	DebugMessages []string
	ErrorMessages []string
}

func NewMockLogger() *MockLogger {
	return &MockLogger{}
}

func (m *MockLogger) Debug(msg string, args ...interface{}) {
	m.DebugMessages = append(m.DebugMessages, msg)
}

func (m *MockLogger) Info(msg string, args ...interface{}) {
	m.InfoMessages = append(m.InfoMessages, msg)
}

func (m *MockLogger) Warn(msg string, args ...interface{}) {
	// Implement this if you use it in your tests
}

func (m *MockLogger) Error(msg string, err error, args ...interface{}) {
	m.ErrorMessages = append(m.ErrorMessages, msg)
}

func (m *MockLogger) Fatal(msg string, err error, args ...interface{}) {
	// Implement this if you use it in your tests
}

func (m *MockLogger) WithOperation(operationID string) LoggerInterface {
	// Implement this if you use it in your tests
	return m
}

func (m *MockLogger) IsDebugEnabled() bool {
	// Implement this if you use it in your tests
	return false
}
