package db

import (
	"context"
	"encoding/json"
	"github.com/CloudImpl-Inc/next-coder-sdk/polycode"
)

type Query struct {
	collection *Collection
	filter     string
	args       []any
	limit      int
}

func (q Query) StartFrom(token polycode.PageToken) (polycode.Query, error) {
	//TODO implement me
	panic("implement me")
}

func (q Query) Index(name string) polycode.Query {
	//TODO implement me
	panic("implement me")
}

func (q Query) Project(paths ...string) polycode.Query {
	//TODO implement me
	panic("implement me")
}

func (q Query) ProjectExpr(expr string, args ...interface{}) polycode.Query {
	//TODO implement me
	panic("implement me")
}

func (q Query) Filter(expr string, args ...interface{}) polycode.Query {
	q.filter = expr
	q.args = args
	return q
}

func (q Query) Consistent(on bool) polycode.Query {
	//TODO implement me
	panic("implement me")
}

func (q Query) Limit(limit int) polycode.Query {
	q.limit = limit
	return q
}

func (q Query) SearchLimit(limit int) polycode.Query {
	//TODO implement me
	panic("implement me")
}

func (q Query) RequestLimit(limit int) polycode.Query {
	//TODO implement me
	panic("implement me")
}

func (q Query) Order(order polycode.Order) polycode.Query {
	//TODO implement me
	panic("implement me")
}

func (q Query) One(ctx context.Context, ret interface{}) error {
	r, err := q.collection.db.client.GetItem(q.collection.db.sessionId, q.collection.name, "", q.filter, q.args)
	if err != nil {
		return err
	}
	b, err := json.Marshal(r)
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
	r, err := q.collection.db.client.QueryItems(q.collection.db.sessionId, q.collection.name, q.filter, q.args, q.limit)
	if err != nil {
		println("error queryitem ", err.Error())
		return err
	}

	b, err := json.Marshal(r)
	if err != nil {
		println("error marshal queryitem ", err.Error())
		return err
	}

	println("query all ", string(b))
	err = json.Unmarshal(b, ret)
	if err != nil {
		println("error unmarshal queryitem ", err.Error())
		return err
	}
	return nil
}

func (q Query) AllWithNextToken(ctx context.Context, ret interface{}) (polycode.PageToken, error) {
	//TODO implement me
	panic("implement me")
}

func (q Query) Iter() polycode.PagingIter {
	//TODO implement me
	panic("implement me")
}
