package polycode

import (
	"encoding/base64"
	"fmt"
)

type ReadOnlyFileStore struct {
	client    *ServiceClient
	sessionId string
}

func (d ReadOnlyFileStore) Folder(name string) ReadOnlyFolder {
	return Folder{
		client:    d.client,
		sessionId: d.sessionId,
		name:      name,
	}
}

type FileStore struct {
	client    *ServiceClient
	sessionId string
}

func (d FileStore) Folder(name string) Folder {
	return Folder{
		client:    d.client,
		sessionId: d.sessionId,
		name:      name,
	}
}

type ReadOnlyFolder interface {
	Load(name string) ([]byte, error)
}

type Folder struct {
	client    *ServiceClient
	sessionId string
	name      string
}

func (f Folder) Load(name string) ([]byte, error) {
	req := GetFileRequest{
		Key: f.name + "/" + name,
	}

	res, err := f.client.GetFile(f.sessionId, req)
	if err != nil {
		fmt.Printf("failed to get file: %s\n", err.Error())
		return nil, err
	}

	// Decode the base64 data
	data, err := base64.StdEncoding.DecodeString(res.Content)
	if err != nil {
		fmt.Printf("failed to decode base64: %s\n", err.Error())
		return nil, err
	}

	return data, nil
}

func (f Folder) Save(name string, data []byte) error {
	// Encode the data as base64
	base64Data := base64.StdEncoding.EncodeToString(data)
	req := PutFileRequest{
		Key:     f.name + "/" + name,
		Content: base64Data,
	}

	err := f.client.PutFile(f.sessionId, req)
	if err != nil {
		fmt.Printf("failed to put file: %s\n", err.Error())
		return err
	}

	return nil
}

func NewFileStore(client *ServiceClient, sessionId string) FileStore {
	return FileStore{
		client:    client,
		sessionId: sessionId,
	}
}

func NewReadOnlyFileStore(client *ServiceClient, sessionId string) ReadOnlyFileStore {
	return ReadOnlyFileStore{
		client:    client,
		sessionId: sessionId,
	}
}
