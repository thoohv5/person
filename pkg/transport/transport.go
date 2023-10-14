package transport

import (
	"context"
	"net/url"
)

// Server is transport server.
type Server interface {
	Start(context.Context) error
	Stop(context.Context) error
}

// Endpointer is registry endpoint.
type Endpointer interface {
	Endpoint() (*url.URL, error)
}

// Header is the storage medium used by a Header.
type Header interface {
	Get(key string) string
	Set(key string, value string)
	Keys() []string
}

// Transporter is transport context value interface.
type Transporter interface {
	// Kind transporter
	// grpc
	// http
	Kind() Kind
	// Endpoint return server or client endpoint
	// Server Transport: grpc://127.0.0.1:9000
	// Client Transport: discovery:///provider-demo
	Endpoint() string
	// Operation Service full method selector generated by protobuf
	// example: /helloworld.Greeter/SayHello
	Operation() string
	// RequestHeader return transport request header
	// http: http.Header
	// grpc: metadata.MD
	RequestHeader() Header
	// ReplyHeader return transport reply/response header
	// only valid for server transport
	// http: http.Header
	// grpc: metadata.MD
	ReplyHeader() Header
}

// Kind defines the type of Transport
type Kind string

func (k Kind) String() string { return string(k) }

// Defines a set of transport kind
const (
	KindGRPC Kind = "grpc"
	KindHTTP Kind = "http"
)
