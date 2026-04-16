package adapters

import (
	"log/slog"
	"os"
)

type OSFileRemover struct {
	logger *slog.Logger
}

func NewOsFileRemover(logger *slog.Logger) OSFileRemover {
	return OSFileRemover{logger}
}

func (fr OSFileRemover) Remove(name string) error {
	fr.logger.Info("Remove ", "file", name)
	return os.Remove(name)
}
