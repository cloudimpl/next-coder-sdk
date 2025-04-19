package polycode

import (
	"context"
	"fmt"
	"log"
)

// Order is used for specifying the order of results.
type Order bool

// Orders for sorting results.
const (
	Ascending  Order = true  // ScanIndexForward = true
	Descending       = false // ScanIndexForward = false
)

type PageToken string

type Iter interface {
	Next(ctx context.Context, out interface{}) bool
	Err() error
}

type PagingIter interface {
	Iter
	NextToken(context.Context) (PageToken, error)
}

type Query struct {
	collection *Collection
	filter     string
	args       []any
	limit      int
}

func (q Query) StartFrom(token PageToken) (Query, error) {
	//TODO implement me
	panic("implement me")
}

func (q Query) Index(name string) Query {
	//TODO implement me
	panic("implement me")
}

func (q Query) Project(paths ...string) Query {
	//TODO implement me
	panic("implement me")
}

func (q Query) ProjectExpr(expr string, args ...interface{}) Query {
	//TODO implement me
	panic("implement me")
}

func (q Query) Filter(expr string, args ...interface{}) Query {
	q.filter = expr
	q.args = args
	return q
}

func (q Query) Consistent(on bool) Query {
	//TODO implement me
	panic("implement me")
}

func (q Query) Limit(limit int) Query {
	q.limit = limit
	return q
}

func (q Query) SearchLimit(limit int) Query {
	//TODO implement me
	panic("implement me")
}

func (q Query) RequestLimit(limit int) Query {
	//TODO implement me
	panic("implement me")
}

func (q Query) Order(order Order) Query {
	//TODO implement me
	panic("implement me")
}

func (q Query) One(ctx context.Context, ret interface{}) (bool, error) {
	req := QueryRequest{
		Collection: q.collection.name,
		Key:        "",
		Filter:     q.filter,
		Args:       q.args,
	}

	r, err := q.collection.client.QueryItems(q.collection.sessionId, req)
	if err != nil {
		fmt.Printf("client: error query item %s\n", err.Error())
		return false, err
	}

	if len(r) == 0 {
		return false, nil
	}

	e := r[0]
	err = ConvertType(e, ret)
	if err != nil {
		fmt.Printf("failed to convert type: %s\n", err.Error())
		return false, err
	}

	return true, nil
}

func (q Query) Count(ctx context.Context) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (q Query) All(ctx context.Context, ret interface{}) error {
	req := QueryRequest{
		Collection: q.collection.name,
		Key:        "",
		Filter:     q.filter,
		Args:       q.args,
		Limit:      q.limit,
	}

	r, err := q.collection.client.QueryItems(q.collection.sessionId, req)
	if err != nil {
		log.Println("client: error query item ", err.Error())
		return err
	}

	err = ConvertType(r, ret)
	if err != nil {
		fmt.Printf("failed to convert type: %s\n", err.Error())
		return err
	}

	return nil
}

func (q Query) AllWithNextToken(ctx context.Context, ret interface{}) (PageToken, error) {
	//TODO implement me
	panic("implement me")
}

func (q Query) Iter() PagingIter {
	//TODO implement me
	panic("implement me")
}

type UnsafeQuery struct {
	tenantId     string
	partitionKey string
	collection   *UnsafeCollection
	filter       string
	args         []any
	limit        int
}

func (q UnsafeQuery) StartFrom(token PageToken) (UnsafeQuery, error) {
	//TODO implement me
	panic("implement me")
}

func (q UnsafeQuery) Index(name string) UnsafeQuery {
	//TODO implement me
	panic("implement me")
}

func (q UnsafeQuery) Project(paths ...string) UnsafeQuery {
	//TODO implement me
	panic("implement me")
}

func (q UnsafeQuery) ProjectExpr(expr string, args ...interface{}) UnsafeQuery {
	//TODO implement me
	panic("implement me")
}

func (q UnsafeQuery) Filter(expr string, args ...interface{}) UnsafeQuery {
	q.filter = expr
	q.args = args
	return q
}

func (q UnsafeQuery) Consistent(on bool) UnsafeQuery {
	//TODO implement me
	panic("implement me")
}

func (q UnsafeQuery) Limit(limit int) UnsafeQuery {
	q.limit = limit
	return q
}

func (q UnsafeQuery) SearchLimit(limit int) UnsafeQuery {
	//TODO implement me
	panic("implement me")
}

func (q UnsafeQuery) RequestLimit(limit int) UnsafeQuery {
	//TODO implement me
	panic("implement me")
}

func (q UnsafeQuery) Order(order Order) UnsafeQuery {
	//TODO implement me
	panic("implement me")
}

func (q UnsafeQuery) One(ctx context.Context, ret interface{}) (bool, error) {
	req := UnsafeQueryRequest{
		TenantId:     q.tenantId,
		PartitionKey: q.partitionKey,
		Collection:   q.collection.name,
		Key:          "",
		Filter:       q.filter,
		Args:         q.args,
	}

	r, err := q.collection.client.UnsafeQueryItems(q.collection.sessionId, req)
	if err != nil {
		fmt.Printf("client: error query item %s\n", err.Error())
		return false, err
	}

	if len(r) == 0 {
		return false, nil
	}

	e := r[0]
	err = ConvertType(e, ret)
	if err != nil {
		fmt.Printf("failed to convert type: %s\n", err.Error())
		return false, err
	}

	return true, nil
}

func (q UnsafeQuery) Count(ctx context.Context) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (q UnsafeQuery) All(ctx context.Context, ret interface{}) error {
	req := UnsafeQueryRequest{
		TenantId:     q.tenantId,
		PartitionKey: q.partitionKey,
		Collection:   q.collection.name,
		Key:          "",
		Filter:       q.filter,
		Args:         q.args,
		Limit:        q.limit,
	}

	r, err := q.collection.client.UnsafeQueryItems(q.collection.sessionId, req)
	if err != nil {
		log.Println("client: error query item ", err.Error())
		return err
	}

	err = ConvertType(r, ret)
	if err != nil {
		fmt.Printf("failed to convert type: %s\n", err.Error())
		return err
	}

	return nil
}

func (q UnsafeQuery) AllWithNextToken(ctx context.Context, ret interface{}) (PageToken, error) {
	//TODO implement me
	panic("implement me")
}

func (q UnsafeQuery) Iter() PagingIter {
	//TODO implement me
	panic("implement me")
}
