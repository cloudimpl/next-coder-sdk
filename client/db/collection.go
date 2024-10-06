package db

import (
	"encoding/json"
	"fmt"
	"github.com/CloudImpl-Inc/next-coder-sdk/polycode"
	"reflect"
)

type Collection struct {
	db   *Database
	dbTx *Tx
	name string
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

	c.dbTx.Operations = append(c.dbTx.Operations, Operation{
		Action:     "insert",
		Collection: c.name,
		Key:        id,
		Item:       mapItems,
	})
	return nil
	//b, err := json.Marshal(item)
	//if err != nil {
	//	return fmt.Errorf("failed to marshal JSON: %w", err)
	//}
	//
	//var mapItems map[string]interface{}
	//err = json.Unmarshal(b, &mapItems)
	//if err != nil {
	//	return fmt.Errorf("failed to unmarshal JSON: %w", err)
	//}
	//return c.db.client.InsertItem(c.db.sessionId, c.name, id, mapItems)
}

func (c Collection) DeleteOne(key string) error {
	c.dbTx.Operations = append(c.dbTx.Operations, Operation{
		Action:     "delete",
		Collection: c.name,
		Key:        key,
	})
	return nil
	//return c.db.client.DeleteItem(c.db.sessionId, c.name, key)
}

func (c Collection) GetOne(key string, ret interface{}) (bool, error) {
	r, err := c.db.client.GetItem(c.db.sessionId, c.name, key, "", nil)
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

func (c Collection) Query() polycode.Query {
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
