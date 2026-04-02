package mocks

import (
	"context"
	"meme-bot/internal/domain"

	"github.com/stretchr/testify/mock"
)

type SourceMock struct {
	mock.Mock
}

func (m *SourceMock) FetchMeme(ctx context.Context) (*domain.Meme, string, error) {
	args := m.Called(ctx)

	var meme *domain.Meme
	if v := args.Get(0); v != nil {
		meme = v.(*domain.Meme)
	}

	return meme, args.String(1), args.Error(2)
}
