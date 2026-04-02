package usecase_test

import (
	"fmt"
	"io"
	"log/slog"
	"meme-bot/internal/domain"
	"meme-bot/internal/mocks"
	"meme-bot/internal/ports"
	usecase "meme-bot/internal/use_case"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPublishSuccess(t *testing.T) {
	repo := new(mocks.RepositoryMock)
	remover := new(mocks.RemoverMock)
	pub1 := new(mocks.PublisherMock)
	pub2 := new(mocks.PublisherMock)

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	uc := usecase.NewPublisherUseCase(
		[]ports.Publisher{pub1, pub2},
		repo,
		logger,
		remover,
	)

	meme := domain.Meme{ID: 1, SourceID: "abc"}

	repo.On("GetOldestApproved").Return(meme, nil)

	pub1.On("Publish", meme).Return(nil)
	pub2.On("Publish", meme).Return(nil)

	repo.On("Delete", 1).Return(nil)
	remover.On("Remove", "tmp/abc.jpg").Return(nil)
	repo.On("UpdateStatus", 1, domain.Posted).Return(nil)

	err := uc.Call()

	assert.NoError(t, err)
}

func TestPublishPartialFailure(t *testing.T) {
	repo := new(mocks.RepositoryMock)
	remover := new(mocks.RemoverMock)
	pub1 := new(mocks.PublisherMock)
	pub2 := new(mocks.PublisherMock)

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	uc := usecase.NewPublisherUseCase(
		[]ports.Publisher{pub1, pub2},
		repo,
		logger,
		remover,
	)

	meme := domain.Meme{ID: 1, SourceID: "abc"}

	repo.On("GetOldestApproved").Return(meme, nil)

	pub1.On("Publish", meme).Return(fmt.Errorf("fail"))
	pub2.On("Publish", meme).Return(nil)

	repo.On("Delete", 1).Return(nil)
	remover.On("Remove", "tmp/abc.jpg").Return(nil)
	repo.On("UpdateStatus", 1, domain.Posted).Return(nil)

	err := uc.Call()

	assert.NoError(t, err)
}

func TestPublishAllFailed(t *testing.T) {
	repo := new(mocks.RepositoryMock)
	remover := new(mocks.RemoverMock)
	pub1 := new(mocks.PublisherMock)
	pub2 := new(mocks.PublisherMock)

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	uc := usecase.NewPublisherUseCase(
		[]ports.Publisher{pub1, pub2},
		repo,
		logger,
		remover,
	)

	meme := domain.Meme{ID: 1, SourceID: "abc"}

	repo.On("GetOldestApproved").Return(meme, nil)

	pub1.On("Publish", meme).Return(fmt.Errorf("fail1"))
	pub2.On("Publish", meme).Return(fmt.Errorf("fail2"))

	repo.On("Delete", 1).Return(nil)
	remover.On("Remove", "tmp/abc.jpg").Return(nil)

	err := uc.Call()

	assert.Error(t, err)
	repo.AssertNotCalled(t, "UpdateStatus", mock.Anything)
}

func TestPublishRemoveError(t *testing.T) {
	repo := new(mocks.RepositoryMock)
	remover := new(mocks.RemoverMock)
	pub := new(mocks.PublisherMock)

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	uc := usecase.NewPublisherUseCase(
		[]ports.Publisher{pub},
		repo,
		logger,
		remover,
	)

	meme := domain.Meme{ID: 1, SourceID: "abc"}

	repo.On("GetOldestApproved").Return(meme, nil)
	pub.On("Publish", meme).Return(nil)

	repo.On("Delete", 1).Return(nil)
	remover.On("Remove", "tmp/abc.jpg").Return(fmt.Errorf("fs error"))

	err := uc.Call()

	assert.Error(t, err)
}
