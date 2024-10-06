package db

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
)

// Client struct to hold the base URL and HTTP client
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// QueryRequest represents the JSON structure for query operations
type QueryRequest struct {
	Collection string        `json:"collection"`
	SessionId  string        `json:"sessionId"`
	Key        string        `json:"key"`
	Filter     string        `json:"filter"`
	Args       []interface{} `json:"args"`
	Limit      int           `json:"limit"`
}

type GetFileRequest struct {
	SessionId string `json:"sessionId"`
	Key       string `json:"key"`
}

type GetFileResponse struct {
	Content string `json:"content"`
}

type PutFileRequest struct {
	SessionId string `json:"sessionId"`
	Key       string `json:"key"`
	Content   string `json:"content"`
}

// DBRequest represents the JSON structure for insert operations
//type DBRequest struct {
//	Collection string                 `json:"collection"`
//	Key        string                 `json:"key"`
//	SessionId  string                 `json:"sessionId"`
//	Obj        map[string]interface{} `json:"Obj"`
//}

// CommitRequest represents a commit operation
//type CommitRequest struct {
//	SessionId string `json:"sessionId"`
//}

// RollbackRequest represents a rollback operation
//type RollbackRequest struct {
//	SessionId string `json:"sessionId"`
//}

// NewClient creates a new instance of Client with the given base URL
func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
	}
}

func (c *Client) GetItem(sessionId string, collection, key string, filter string, args []any) (map[string]interface{}, error) {
	reqBody, err := json.Marshal(QueryRequest{
		Collection: collection,
		SessionId:  sessionId,
		Key:        key,
		Args:       args,
		Filter:     filter,
	})
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Post(fmt.Sprintf("%s/db/get", c.BaseURL), "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		println("GetItem", err.Error())
		return nil, err
	}

	println(fmt.Sprintf("db get item received "))
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get item, status: %v", resp.Status)
	}

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	println(fmt.Sprintf("db get item received response %v", result))
	return result, nil
}

// QueryItems items with filters
func (c *Client) QueryItems(sessionId, collection, filter string, args []any, limit int) ([]map[string]interface{}, error) {
	reqBody, err := json.Marshal(QueryRequest{
		Collection: collection,
		SessionId:  sessionId,
		Filter:     filter,
		Args:       args,
		Limit:      limit,
	})
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Post(fmt.Sprintf("%s/db/query", c.BaseURL), "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		println("QueryItems", err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	var result []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	// Print out the result
	for _, item := range result {
		fmt.Printf("Obj: %+v\n", item)
	}

	return result, nil
}

func (c *Client) GetFile(sessionId string, key string) ([]byte, error) {
	reqBody, err := json.Marshal(GetFileRequest{
		SessionId: sessionId,
		Key:       key,
	})
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Post(fmt.Sprintf("%s/file/get", c.BaseURL), "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		println("GetFile", err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get file, status: %v", resp.Status)
	}

	var result GetFileResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	// Decode the Base64 string
	decodedData, err := base64.StdEncoding.DecodeString(result.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 data: %v", err)
	}

	return decodedData, nil
}

func (c *Client) PutFile(sessionId string, key string, content []byte) error {
	// Encode the data as base64
	base64Data := base64.StdEncoding.EncodeToString(content)
	reqBody, err := json.Marshal(PutFileRequest{
		SessionId: sessionId,
		Key:       key,
		Content:   base64Data,
	})
	if err != nil {
		return err
	}

	resp, err := c.HTTPClient.Post(fmt.Sprintf("%s/file/put", c.BaseURL), "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		println("PutFile", err.Error())
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to put file, status: %v", resp.Status)
	}

	return nil
}

// InsertItem Insert an item into the collection
//func (c *Client) InsertItem(sessionId string, collection, key string, item map[string]interface{}) error {
//	reqBody, err := json.Marshal(DBRequest{
//		Collection: collection,
//		SessionId:  sessionId,
//		Key:        key,
//		Obj:        item,
//	})
//	if err != nil {
//		return err
//	}
//
//	resp, err := c.HTTPClient.Post(fmt.Sprintf("%s/insert", c.BaseURL), "application/json", bytes.NewBuffer(reqBody))
//	if err != nil {
//		println("InsertItem", err.Error())
//		return err
//	}
//	defer resp.Body.Close()
//
//	if resp.StatusCode != http.StatusOK {
//		return fmt.Errorf("failed to insert item, status: %v", resp.Status)
//	}
//	return nil
//}

//func (c *Client) DeleteItem(sessionId string, collection, key string) error {
//	reqBody, err := json.Marshal(DBRequest{
//		Collection: collection,
//		SessionId:  sessionId,
//		Key:        key,
//	})
//	if err != nil {
//		return err
//	}
//
//	resp, err := c.HTTPClient.Post(fmt.Sprintf("%s/delete", c.BaseURL), "application/json", bytes.NewBuffer(reqBody))
//	if err != nil {
//		println("DeleteItem", err.Error())
//		return err
//	}
//	defer resp.Body.Close()
//
//	if resp.StatusCode != http.StatusOK {
//		return fmt.Errorf("failed to delete item, status: %v", resp.Status)
//	}
//	return nil
//}

// Commit the transaction
//func (c *Client) CommitTransaction(sessionId string) error {
//	reqBody, _ := json.Marshal(CommitRequest{
//		SessionId: sessionId,
//	})
//
//	resp, err := c.HTTPClient.Post(fmt.Sprintf("%s/commit", c.BaseURL), "application/json", bytes.NewBuffer(reqBody))
//	if err != nil {
//		println("CommitTransaction", err.Error())
//		return err
//	}
//	defer resp.Body.Close()
//
//	if resp.StatusCode != http.StatusOK {
//		return fmt.Errorf("failed to commit transaction, status: %v", resp.Status)
//	}
//	return nil
//}

// Rollback the transaction
//func (c *Client) RollbackTransaction(sessionId string) error {
//	reqBody, _ := json.Marshal(RollbackRequest{
//		SessionId: sessionId,
//	})
//
//	resp, err := c.HTTPClient.Post(fmt.Sprintf("%s/rollback", c.BaseURL), "application/json", bytes.NewBuffer(reqBody))
//	if err != nil {
//		return err
//	}
//	defer resp.Body.Close()
//
//	if resp.StatusCode != http.StatusOK {
//		return fmt.Errorf("failed to rollback transaction, status: %v", resp.Status)
//	}
//	return nil
//}

//func main() {
//	// Initialize the client with the server base URL
//	client := NewClient("http://localhost:8080")
//
//	// Example usage of client operations
//	item := map[string]interface{}{
//		"Name": "John Doe",
//		"Age":  30,
//	}
//
//	err := client.InsertItem("1", "users", "user123", item)
//	if err != nil {
//		fmt.Println("Error inserting item:", err)
//		return
//	}
//
//	_, err = client.QueryItems("1", "users", "Age >= 25 AND Name == \"John Doe\"", nil, 10)
//	if err != nil {
//		fmt.Println("Error querying items:", err)
//		return
//	}
//
//	err = client.CommitTransaction("1")
//	if err != nil {
//		fmt.Println("Error committing transaction:", err)
//		return
//	}
//}
