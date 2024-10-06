package db

import "github.com/CloudImpl-Inc/next-coder-sdk/polycode"

type Tx struct {
	Operations []Operation `json:"operations"`
}

type Operation struct {
	Action     string                 `json:"action"`
	Collection string                 `json:"collection"`
	Key        string                 `json:"key"`
	Item       map[string]interface{} `json:"item"`
}

type Database struct {
	client    *Client
	sessionId string
	tx        *Tx
}

func (d Database) GetTx() *Tx {
	return d.tx
}

func (d Database) Collection(name string) polycode.Collection {
	return Collection{
		db:   &d,
		dbTx: d.tx,
		name: name,
	}
}

func NewDatabase(client *Client, sessionId string) Database {
	return Database{
		client:    client,
		sessionId: sessionId,
		tx:        &Tx{},
	}
}
