package polycode

import (
	"fmt"
	"reflect"
	"time"
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
	return c.InsertOneWithTTL(item, -1)
}

func (c Collection) InsertOneWithTTL(item interface{}, expireIn time.Duration) error {
	var ttl int64
	if expireIn == -1 {
		ttl = -1
	} else {
		ttl = time.Now().Unix() + int64(expireIn.Seconds())
	}

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
		TTL:        ttl,
	}

	err = c.client.PutItem(c.sessionId, req)
	if err != nil {
		fmt.Printf("failed to put item: %s\n", err.Error())
		return err
	}

	return nil
}

func (c Collection) UpdateOne(item interface{}) error {
	return c.UpdateOneWithTTL(item, -1)
}

func (c Collection) UpdateOneWithTTL(item interface{}, expireIn time.Duration) error {
	var ttl int64
	if expireIn == -1 {
		ttl = -1
	} else {
		ttl = time.Now().Unix() + int64(expireIn.Seconds())
	}

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
		TTL:        ttl,
	}

	err = c.client.PutItem(c.sessionId, req)
	if err != nil {
		fmt.Printf("failed to put item: %s\n", err.Error())
		return err
	}

	return nil
}

func (c Collection) UpsertOne(item interface{}) error {
	return c.UpsertOneWithTTL(item, -1)
}

func (c Collection) UpsertOneWithTTL(item interface{}, expireIn time.Duration) error {
	var ttl int64
	if expireIn == -1 {
		ttl = -1
	} else {
		ttl = time.Now().Unix() + int64(expireIn.Seconds())
	}

	id, err := GetId(item)
	if err != nil {
		fmt.Printf("failed to get id: %s\n", err.Error())
		return err
	}

	req := PutRequest{
		Action:     "upsert",
		Collection: c.name,
		Key:        id,
		Item:       item,
		TTL:        ttl,
	}

	err = c.client.PutItem(c.sessionId, req)
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

	err := c.client.PutItem(c.sessionId, req)
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

	r, err := c.client.GetItem(c.sessionId, req)
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

func newDatabase(client *ServiceClient, sessionId string) DataStore {
	return DataStore{
		client:    client,
		sessionId: sessionId,
	}
}
