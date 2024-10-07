package loggo

import (
	"github.com/golang/mock/gomock"
)

// NewMockLogger creates a new mock logger for testing purposes
func NewMockLogger(ctrl *gomock.Controller) *MockLoggerInterface {
	return NewMockLoggerInterface(ctrl)
}
