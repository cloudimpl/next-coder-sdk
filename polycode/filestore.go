package polycode

type FileStore interface {
	Folder(name string) Folder
}

type ReadOnlyFileStore interface {
	Folder(name string) ReadOnlyFolder
}

type ReadOnlyFolder interface {
	Load(name string) ([]byte, error)
}

type Folder interface {
	Load(name string) ([]byte, error)
	Save(name string, data []byte) error
}
