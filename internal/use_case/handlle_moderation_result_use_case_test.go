package usecase_test

import (
	"errors"
	"meme-bot/internal/domain"
	"meme-bot/internal/mocks"
	usecase "meme-bot/internal/use_case"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleModerationApproved(t *testing.T) {
	repo := new(mocks.RepositoryMock)
	remover := new(mocks.RemoverMock)

	uc := usecase.NewHandleModerationResultUseCase(repo, remover)

	repo.
		On("UpdateStatus", 1, domain.Approved).
		Return(nil)

	err := uc.Call(1, usecase.Approved)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestHandleModerationRejectedSuccess(t *testing.T) {
	repo := new(mocks.RepositoryMock)
	remover := new(mocks.RemoverMock)

	uc := usecase.NewHandleModerationResultUseCase(repo, remover)

	meme := domain.Meme{
		ID:       1,
		SourceID: "abc",
	}

	repo.On("GetByID", 1).Return(meme, nil)
	repo.On("Delete", 1).Return(nil)
	remover.On("Remove", "tmp/abc.jpg").Return(nil)

	err := uc.Call(1, usecase.Rejected)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
	remover.AssertExpectations(t)
}

func TestHandleModerationRejectedGetError(t *testing.T) {
	repo := new(mocks.RepositoryMock)
	remover := new(mocks.RemoverMock)

	uc := usecase.NewHandleModerationResultUseCase(repo, remover)

	repo.
		On("GetByID", 1).
		Return(domain.Meme{}, errors.New("db error"))

	err := uc.Call(1, usecase.Rejected)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed on getting meme by ID")
}

func TestHandleModerationRejectedRemoveError(t *testing.T) {
	repo := new(mocks.RepositoryMock)
	remover := new(mocks.RemoverMock)

	uc := usecase.NewHandleModerationResultUseCase(repo, remover)

	meme := domain.Meme{
		ID:       1,
		SourceID: "abc",
	}

	repo.On("GetByID", 1).Return(meme, nil)
	repo.On("Delete", 1).Return(nil)
	remover.On("Remove", "tmp/abc.jpg").Return(errors.New("fs error"))

	err := uc.Call(1, usecase.Rejected)

	assert.Error(t, err)
}

func TestHandleModerationInvalidSelection(t *testing.T) {
	repo := new(mocks.RepositoryMock)
	remover := new(mocks.RemoverMock)

	uc := usecase.NewHandleModerationResultUseCase(repo, remover)

	err := uc.Call(1, "unknown")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no such user selection")
}
