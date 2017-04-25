package loadbalancer

import (
	"github.com/skynetservices/skynet"
	"errors"
)

var (
	NoInstances = errors.New("No instance")
)

type LoadBalancer interface {
	AddInstance(s skynet.ServiceInfo)
	UpdateInstance(s skynet.ServiceInfo)
	RemoveInstance(s skynet.ServiceInfo)
	Choose()(skynet.ServiceInfo, error)
}
