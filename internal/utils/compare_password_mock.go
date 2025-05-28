package utils

import "github.com/stretchr/testify/mock"

type ComparePassMock struct {
	mock.Mock
}

func NewComparePassMock() *ComparePassMock {
	return &ComparePassMock{}
}

func (m *ComparePassMock) ComparePassword(hashedPassword, inputPassword string) error {
	args := m.Called(hashedPassword,inputPassword)
	return args.Error(0)
}