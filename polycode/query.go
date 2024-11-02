package polycode

import (
	"context"
	"encoding/json"
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

func (q Query) One(ctx context.Context, ret interface{}) error {
	req := QueryRequest{
		Collection: q.collection.name,
		Key:        "",
		Filter:     q.filter,
		Args:       q.args,
	}

	r, err := q.collection.client.QueryItems(q.collection.sessionId, req)
	if err != nil {
		return err
	}

	if len(r) == 0 {
		return nil
	}

	e := r[0]
	b, err := json.Marshal(e)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, ret)
	if err != nil {
		return err
	}

	return nil
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

	b, err := json.Marshal(r)
	if err != nil {
		log.Println("client: error marshal query item ", err.Error())
		return err
	}

	err = json.Unmarshal(b, ret)
	if err != nil {
		log.Println("client: error unmarshal query item ", err.Error())
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
