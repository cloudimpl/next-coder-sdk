package polycode

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type DataStore struct {
	client    *ServiceClient
	sessionId string
}

func (d DataStore) Collection(name string) Collection {
	return Collection{
		client:    d.client,
		sessionId: d.sessionId,
		name:      name,
	}
}

type Collection struct {
	client    *ServiceClient
	sessionId string
	name      string
}

func (c Collection) InsertOne(item interface{}) error {
	id, err := GetId(item)
	if err != nil {
		return err
	}

	b, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	var mapItems map[string]interface{}
	err = json.Unmarshal(b, &mapItems)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	req := PutRequest{
		Action:     "insert",
		Collection: c.name,
		Key:        id,
		Item:       mapItems,
	}

	err = c.client.PutItem(c.sessionId, req)
	if err != nil {
		return err
	}

	return nil
}

func (c Collection) UpdateOne(item interface{}) error {
	id, err := GetId(item)
	if err != nil {
		return err
	}

	b, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	var mapItems map[string]interface{}
	err = json.Unmarshal(b, &mapItems)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	req := PutRequest{
		Action:     "update",
		Collection: c.name,
		Key:        id,
		Item:       mapItems,
	}

	err = c.client.PutItem(c.sessionId, req)
	if err != nil {
		return err
	}

	return nil
}

func (c Collection) DeleteOne(key string) error {
	req := PutRequest{
		Action:     "delete",
		Collection: c.name,
		Key:        key,
	}

	err := c.client.PutItem(c.sessionId, req)
	if err != nil {
		return err
	}

	return nil
}

func (c Collection) GetOne(key string, ret interface{}) (bool, error) {
	req := QueryRequest{
		Collection: c.name,
		Key:        key,
		Filter:     "",
		Args:       nil,
	}

	r, err := c.client.GetItem(c.sessionId, req)
	if err != nil {
		return false, err
	}

	if r == nil {
		return false, nil
	}

	b, err := json.Marshal(r)
	if err != nil {
		return false, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	err = json.Unmarshal(b, ret)
	if err != nil {
		return false, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return true, nil
}

func (c Collection) Query() Query {
	return Query{
		collection: &c,
	}
}

func GetId(item any) (string, error) {
	id := ""
	v := reflect.ValueOf(item)
	t := reflect.TypeOf(item)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i).Interface()

		// Skip the PKEY and RKEY fields
		if field.Tag.Get("polycode") == "id" {
			id = value.(string)
			break
		}
	}

	if id == "" {
		return "", fmt.Errorf("id not found")
	}
	return id, nil
}

func NewDatabase(client *ServiceClient, sessionId string) DataStore {
	return DataStore{
		client:    client,
		sessionId: sessionId,
	}
}
