package polycode

import (
	"encoding/base64"
	client2 "github.com/CloudImpl-Inc/next-coder-sdk/client"
)

type FileStore struct {
	client    *client2.ServiceClient
	sessionId string
}

func (d FileStore) Folder(name string) Folder {
	return Folder{
		client:    d.client,
		sessionId: d.sessionId,
		name:      name,
	}
}

type Folder struct {
	client    *client2.ServiceClient
	sessionId string
	name      string
}

func (f Folder) Load(name string) ([]byte, error) {
	req := client2.GetFileRequest{
		Key: f.name + "/" + name,
	}

	res, err := f.client.GetFile(f.sessionId, req)
	if err != nil {
		return nil, err
	}

	// Decode the base64 data
	data, err := base64.StdEncoding.DecodeString(res.Content)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (f Folder) Save(name string, data []byte) error {
	// Encode the data as base64
	base64Data := base64.StdEncoding.EncodeToString(data)
	req := client2.PutFileRequest{
		Key:     f.name + "/" + name,
		Content: base64Data,
	}

	err := f.client.PutFile(f.sessionId, req)
	if err != nil {
		return err
	}

	return nil
}

func NewFileStore(client *client2.ServiceClient, sessionId string) FileStore {
	return FileStore{
		client:    client,
		sessionId: sessionId,
	}
}
