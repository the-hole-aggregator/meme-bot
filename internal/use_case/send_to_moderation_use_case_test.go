package usecase_test

import (
	"fmt"
	"meme-bot/internal/domain"
	"meme-bot/internal/mocks"
	usecase "meme-bot/internal/use_case"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSendToModerationSuccess(t *testing.T) {
	repo := new(mocks.RepositoryMock)
	pub := new(mocks.PublisherMock)

	uc := usecase.NewSendToModerationUseCase(pub, repo)

	meme := domain.Meme{ID: 1}

	repo.On("GetOldestPending").Return(meme, nil)
	pub.On("Publish", meme).Return(nil)

	err := uc.Call()

	assert.NoError(t, err)
	repo.AssertExpectations(t)
	pub.AssertExpectations(t)
}

func TestSendToModerationRepoError(t *testing.T) {
	repo := new(mocks.RepositoryMock)
	pub := new(mocks.PublisherMock)

	uc := usecase.NewSendToModerationUseCase(pub, repo)

	repo.On("GetOldestPending").Return(domain.Meme{}, fmt.Errorf("db error"))

	err := uc.Call()

	assert.Error(t, err)

	pub.AssertNotCalled(t, "Publish", mock.Anything)
}

func TestSendToModerationPublishError(t *testing.T) {
	repo := new(mocks.RepositoryMock)
	pub := new(mocks.PublisherMock)

	uc := usecase.NewSendToModerationUseCase(pub, repo)

	meme := domain.Meme{ID: 1}

	repo.On("GetOldestPending").Return(meme, nil)
	pub.On("Publish", meme).Return(fmt.Errorf("tg error"))

	err := uc.Call()

	assert.Error(t, err)
}
