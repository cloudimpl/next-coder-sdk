package polycode

import (
	"encoding/base64"
	"errors"
	"fmt"
)

type FileStore struct {
	client    *ServiceClient
	sessionId string
}

func (d FileStore) NewFolder(name string) (Folder, error) {
	req := CreateFolderRequest{
		Folder: name,
	}

	err := d.client.CreateFolder(d.sessionId, req)
	if err != nil {
		fmt.Printf("failed to create folder: %s\n", err.Error())
		return Folder{}, err
	}

	return d.Folder(name), nil
}

func (d FileStore) List(path string, limit int32, nextToken *string) (ListFilePageResponse, error) {
	req := ListFilePageRequest{
		Prefix:            path,
		MaxKeys:           limit,
		ContinuationToken: nextToken,
	}

	res, err := d.client.ListFile(d.sessionId, req)
	if err != nil {
		fmt.Printf("failed to list file: %s\n", err.Error())
		return ListFilePageResponse{}, err
	}

	return res, nil
}

func (d FileStore) Get(path string) (bool, []byte, error) {
	req := GetFileRequest{
		Key: path,
	}

	res, err := d.client.GetFile(d.sessionId, req)
	if err != nil {
		fmt.Printf("failed to get file: %s\n", err.Error())
		return false, nil, err
	}

	if res.Content == "" {
		return false, nil, nil
	}

	// Decode the base64 data
	data, err := base64.StdEncoding.DecodeString(res.Content)
	if err != nil {
		fmt.Printf("failed to decode base64: %s\n", err.Error())
		return true, nil, err
	}

	return true, data, nil
}

func (d FileStore) GetDownloadLink(path string) (string, error) {
	req := GetFileRequest{
		Key: path,
	}

	res, err := d.client.GetFileDownloadLink(d.sessionId, req)
	if err != nil {
		fmt.Printf("failed to get file link: %s\n", err.Error())
		return "", err
	}

	if res.Link == "" {
		return "", errors.New("empty link")
	}

	return res.Link, nil
}

func (d FileStore) Save(path string, data []byte) error {
	// Encode the data as base64
	base64Data := base64.StdEncoding.EncodeToString(data)
	req := PutFileRequest{
		Key:     path,
		Content: base64Data,
	}

	err := d.client.PutFile(d.sessionId, req)
	if err != nil {
		fmt.Printf("failed to put file: %s\n", err.Error())
		return err
	}

	return nil
}

func (d FileStore) GetUploadLink(path string) (string, error) {
	req := GetFileRequest{
		Key: path,
	}

	res, err := d.client.GetFileUploadLink(d.sessionId, req)
	if err != nil {
		fmt.Printf("failed to get file link: %s\n", err.Error())
		return "", err
	}

	if res.Link == "" {
		return "", errors.New("empty link")
	}

	return res.Link, nil
}

func (d FileStore) Delete(path string) error {
	req := DeleteFileRequest{
		Key: path,
	}

	err := d.client.DeleteFile(d.sessionId, req)
	if err != nil {
		fmt.Printf("failed to delete file: %s\n", err.Error())
		return err
	}

	return nil
}

func (d FileStore) Move(oldPath string, newPath string) error {
	req := RenameFileRequest{
		OldKey: oldPath,
		NewKey: newPath,
	}

	err := d.client.RenameFile(d.sessionId, req)
	if err != nil {
		fmt.Printf("failed to rename file: %s\n", err.Error())
		return err
	}

	return nil
}

func (d FileStore) Folder(name string) Folder {
	return Folder{
		client:    d.client,
		sessionId: d.sessionId,
		name:      name,
	}
}

type Folder struct {
	client    *ServiceClient
	sessionId string
	name      string
}

func (f Folder) Load(name string) (bool, []byte, error) {
	req := GetFileRequest{
		Key: f.name + "/" + name,
	}

	res, err := f.client.GetFile(f.sessionId, req)
	if err != nil {
		fmt.Printf("failed to get file: %s\n", err.Error())
		return false, nil, err
	}

	if res.Content == "" {
		return false, nil, nil
	}

	// Decode the base64 data
	data, err := base64.StdEncoding.DecodeString(res.Content)
	if err != nil {
		fmt.Printf("failed to decode base64: %s\n", err.Error())
		return true, nil, err
	}

	return true, data, nil
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

func newFileStore(client *ServiceClient, sessionId string) FileStore {
	return FileStore{
		client:    client,
		sessionId: sessionId,
	}
}
