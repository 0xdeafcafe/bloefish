package contexts

import (
	"context"

	"github.com/0xdeafcafe/bloefish/libraries/config"
)

// ServiceInfo type holds useful info about the currently-running service
type ServiceInfo struct {
	Environment     config.Environment
	Service         string
	ServiceHTTPName string
	GitRepository   string
}

type serviceContextKey string

var serviceContextInfoKey = serviceContextKey("service_context:info")

// SetServiceInfo wraps the context with the service info
func SetServiceInfo(ctx context.Context, serviceInfo ServiceInfo) context.Context {
	return context.WithValue(ctx, serviceContextInfoKey, serviceInfo)
}

// GetServiceInfo retrieves the service info from the context
func GetServiceInfo(ctx context.Context) *ServiceInfo {
	if val, ok := ctx.Value(serviceContextInfoKey).(ServiceInfo); ok {
		return &val
	}

	return nil
}

// MustGetServiceInfo retrieves the service info from the context
func MustGetServiceInfo(ctx context.Context) *ServiceInfo {
	if val, ok := ctx.Value(serviceContextInfoKey).(ServiceInfo); ok {
		return &val
	}

	panic("service info not found in context")
}
