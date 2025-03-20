package polycode

type Counter struct {
	client    *ServiceClient
	sessionId string
	group     string
	name      string
	ttl       int64
}

func (c *Counter) Get() (uint64, error) {
	panic("not implemented")
}

func (c *Counter) Increment(count uint64) (uint64, bool, error) {
	return c.IncrementWithLimit(count, -1)
}

func (c *Counter) IncrementWithLimit(count uint64, limit uint64) (uint64, bool, error) {
	req := IncrementCounterRequest{
		Group: c.group,
		Name:  c.name,
		Count: count,
		Limit: limit,
		TTL:   c.ttl,
	}

	res, err := c.client.IncrementCounter(c.sessionId, req)
	if err != nil {
		return 0, false, err
	}

	return res.Value, res.Incremented, nil
}

func (c *Counter) Decrement(count uint64) (uint64, bool, error) {
	panic("not implemented")
}
