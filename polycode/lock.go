package polycode

import "time"

type Lock struct {
	client    *ServiceClient
	sessionId string
	key       string
}

func (l *Lock) Acquire(expireIn time.Duration) error {
	ttl := time.Now().Unix() + int64(expireIn.Seconds())

	req := AcquireLockRequest{
		Key: l.key,
		TTL: ttl,
	}

	return l.client.AcquireLock(l.sessionId, req)
}

func (l *Lock) Release() error {
	req := ReleaseLockRequest{
		Key: l.key,
	}

	return l.client.ReleaseLock(l.sessionId, req)
}
