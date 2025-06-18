package polycode

type RealtimeChannel struct {
	name          string
	sessionId     string
	serviceClient *ServiceClient
}

func (r RealtimeChannel) Emit(data any) error {
	req := RealtimeEventEmitRequest{
		Channel: r.name,
		Input:   data,
	}

	return r.serviceClient.EmitRealtimeEvent(r.sessionId, req)
}
