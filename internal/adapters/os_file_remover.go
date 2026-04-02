package adapters

import "os"

type OSFileRemover struct{}

func (OSFileRemover) Remove(name string) error {
	return os.Remove(name)
}
