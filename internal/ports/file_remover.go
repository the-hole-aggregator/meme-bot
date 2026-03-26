package ports

type FileRemover interface {
	Remove(name string) error
}
