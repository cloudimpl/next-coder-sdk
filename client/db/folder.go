package db

type Folder struct {
	fileStore *FileStore
	name      string
}

func (f Folder) Load(name string) ([]byte, error) {
	return f.fileStore.client.GetFile(f.fileStore.sessionId, f.name+"/"+name)
}

func (f Folder) Save(name string, data []byte) error {
	return f.fileStore.client.PutFile(f.fileStore.sessionId, f.name+"/"+name, data)
}
