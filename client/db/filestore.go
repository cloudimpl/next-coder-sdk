package db

import "github.com/CloudImpl-Inc/next-coder-sdk/polycode"

type FileStore struct {
	client    *Client
	sessionId string
}

func (d FileStore) Folder(name string) polycode.Folder {
	return Folder{
		fileStore: &d,
		name:      name,
	}
}

func NewFileStore(client *Client, sessionId string) FileStore {
	return FileStore{
		client:    client,
		sessionId: sessionId,
	}
}
