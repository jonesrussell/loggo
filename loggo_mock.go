package loggo

import (
	"github.com/stretchr/testify/mock"
)

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Debug(msg string, args ...interface{}) {
	m.Called(msg, args)
}

func (m *MockLogger) Info(msg string, args ...interface{}) {
	m.Called(msg, args)
}

func (m *MockLogger) Warn(msg string, args ...interface{}) {
	m.Called(msg, args)
}

func (m *MockLogger) Error(msg string, err error, args ...interface{}) {
	m.Called(msg, err, args)
}

func (m *MockLogger) Fatal(msg string, err error, args ...interface{}) {
	m.Called(msg, err, args)
}

func (m *MockLogger) WithOperation(operationID string) LoggerInterface {
	args := m.Called(operationID)
	return args.Get(0).(LoggerInterface)
}

func (m *MockLogger) IsDebugEnabled() bool {
	args := m.Called()
	return args.Bool(0)
}

// Ensure MockLogger implements LoggerInterface
var _ LoggerInterface = &MockLogger{}
