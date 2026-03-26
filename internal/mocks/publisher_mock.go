package mocks

import (
	"meme-bot/internal/domain"

	"github.com/stretchr/testify/mock"
)

type PublisherMock struct {
	mock.Mock
}

func (m *PublisherMock) Publish(meme domain.Meme) error {
	args := m.Called(meme)
	return args.Error(0)
}
