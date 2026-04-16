package usecase_test

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"log/slog"
	"meme-bot/internal/domain"
	"meme-bot/internal/mocks"
	"meme-bot/internal/ports"
	usecase "meme-bot/internal/use_case"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIngestionUseCaseCallSuccess(t *testing.T) {
	ctx := context.Background()
	tempFile := "test_image.png"
	defer func() {
		if err := os.Remove(tempFile); err != nil {
			fmt.Printf("failed on remove file %s", err)
		}
	}()

	assert.NoError(t, createTempImage(tempFile))

	mockRepo := new(mocks.RepositoryMock)
	mockSource := new(mocks.SourceMock)
	mockRemover := new(mocks.RemoverMock)

	meme := &domain.Meme{PHash: "hash123"}

	mockSource.On("FetchMeme", ctx).Return(meme, tempFile, nil).Once()
	mockRepo.On("ExistsByHash", meme.PHash).Return(false, nil).Once()
	mockRepo.On("ExistsBySourceID", meme.SourceID).Return(false, nil).Once()
	mockRepo.On("Save", meme).Return(nil).Once()

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	uc := usecase.NewIngestionUseCase(mockRepo, []ports.Source{mockSource}, logger, mockRemover)
	err := uc.Call(ctx, 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockSource.AssertExpectations(t)
	mockRemover.AssertExpectations(t)
}

func createTempImage(path string) error {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			img.Set(x, y, color.White)
		}
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("failed to close file %s", err)
		}
	}()

	return png.Encode(f, img)
}
