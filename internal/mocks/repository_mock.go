package mocks

import (
	"meme-bot/internal/domain"

	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func (m *RepositoryMock) GetByID(id int) (domain.Meme, error) {
	args := m.Called(id)
	return args.Get(0).(domain.Meme), args.Error(1)
}

func (m *RepositoryMock) Save(meme *domain.Meme) error {
	args := m.Called(meme)
	return args.Error(0)
}

func (m *RepositoryMock) ExistsByHash(hash string) (bool, error) {
	args := m.Called(hash)
	return args.Bool(0), args.Error(1)
}

func (m *RepositoryMock) GetByStatus(status domain.MemeStatus) ([]domain.Meme, error) {
	args := m.Called(status)
	return args.Get(0).([]domain.Meme), args.Error(1)
}

func (m *RepositoryMock) UpdateStatus(id int, status domain.MemeStatus) error {
	args := m.Called(id, status)
	return args.Error(0)
}

func (m *RepositoryMock) GetOldestApproved() (domain.Meme, error) {
	args := m.Called()
	return args.Get(0).(domain.Meme), args.Error(1)
}

func (m *RepositoryMock) GetOldestPending() (domain.Meme, error) {
	args := m.Called()
	return args.Get(0).(domain.Meme), args.Error(1)
}

func (m *RepositoryMock) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
