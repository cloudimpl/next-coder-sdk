package polycode

import (
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
		isGlobal:  false,
	}
}

func (d DataStore) GlobalCollection(name string) Collection {
	return Collection{
		client:    d.client,
		sessionId: d.sessionId,
		name:      name,
		isGlobal:  true,
	}
}

type Collection struct {
	client    *ServiceClient
	sessionId string
	name      string
	isGlobal  bool
}

func (c Collection) InsertOne(item interface{}) error {
	id, err := GetId(item)
	if err != nil {
		fmt.Printf("failed to get id: %s\n", err.Error())
		return err
	}

	req := PutRequest{
		Action:     "insert",
		Collection: c.name,
		Key:        id,
		Item:       item,
	}

	if c.isGlobal {
		err = c.client.PutGlobalItem(c.sessionId, req)
	} else {
		err = c.client.PutItem(c.sessionId, req)
	}
	if err != nil {
		fmt.Printf("failed to put item: %s\n", err.Error())
		return err
	}

	return nil
}

func (c Collection) UpdateOne(item interface{}) error {
	id, err := GetId(item)
	if err != nil {
		fmt.Printf("failed to get id: %s\n", err.Error())
		return err
	}

	req := PutRequest{
		Action:     "update",
		Collection: c.name,
		Key:        id,
		Item:       item,
	}

	if c.isGlobal {
		err = c.client.PutGlobalItem(c.sessionId, req)
	} else {
		err = c.client.PutItem(c.sessionId, req)
	}
	if err != nil {
		fmt.Printf("failed to put item: %s\n", err.Error())
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

	var err error
	if c.isGlobal {
		err = c.client.PutGlobalItem(c.sessionId, req)
	} else {
		err = c.client.PutItem(c.sessionId, req)
	}
	if err != nil {
		fmt.Printf("failed to put item: %s\n", err.Error())
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

	var r map[string]interface{}
	var err error
	if c.isGlobal {
		r, err = c.client.GetGlobalItem(c.sessionId, req)
	} else {
		r, err = c.client.GetItem(c.sessionId, req)
	}
	if err != nil {
		fmt.Printf("failed to get item: %s\n", err.Error())
		return false, err
	}

	if r == nil {
		println("item not found")
		return false, nil
	}

	err = ConvertType(r, ret)
	if err != nil {
		fmt.Printf("failed to convert type: %s\n", err.Error())
		return false, err
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
