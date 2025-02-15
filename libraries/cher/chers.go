package cher

type Transport string

const (
	TransportRPC  Transport = "rpc"
	TransportHTTP Transport = "http"
	TransportWS   Transport = "ws"
)

func TransportNotSupported(requested Transport, supported ...Transport) E {
	return New("transport_not_supported", M{
		"requested": requested,
		"supported": supported,
	})
}
