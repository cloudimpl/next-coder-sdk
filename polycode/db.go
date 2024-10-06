package polycode

import (
	"context"
)

type Database interface {
	Collection(name string) Collection
}

type ReadOnlyDatabase interface {
	Collection(name string) ReadOnlyCollection
}

type ReadOnlyCollection interface {
	GetOne(key string, ret interface{}) error
	Query() Query
}

type Collection interface {
	InsertOne(item interface{}) error
	DeleteOne(key string) error
	GetOne(key string, ret interface{}) (bool, error)
	Query() Query
}

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

type Query interface {
	StartFrom(token PageToken) (Query, error)
	Project(paths ...string) Query
	ProjectExpr(expr string, args ...interface{}) Query
	Filter(expr string, args ...interface{}) Query
	Consistent(on bool) Query
	Limit(limit int) Query
	SearchLimit(limit int) Query
	RequestLimit(limit int) Query
	Order(order Order) Query
	One(ctx context.Context, ret interface{}) error
	Count(ctx context.Context) (int, error)
	All(ctx context.Context, ret interface{}) error
	AllWithNextToken(ctx context.Context, ret interface{}) (PageToken, error)
	Iter() PagingIter
}
