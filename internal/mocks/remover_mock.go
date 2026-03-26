package mocks

import "github.com/stretchr/testify/mock"

type RemoverMock struct {
	mock.Mock
}

func (m *RemoverMock) Remove(path string) error {
	args := m.Called(path)
	return args.Error(0)
}
