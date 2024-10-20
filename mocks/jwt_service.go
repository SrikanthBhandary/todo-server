package mocks

import "github.com/stretchr/testify/mock"

// MockJWTValidator is the mock implementation of JWTValidator
type MockJWTValidator struct {
	mock.Mock
}

// Mock implementation of ValidateToken
func (m *MockJWTValidator) ValidateToken(tokenString string) (int, error) {
	return 1, nil
}

// Mock implementation of GenerateToken
func (m *MockJWTValidator) GenerateToken(userID int) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}
