package structs

import (
	"fmt"
	"reflect"
)

// ConsulConnect represents a Consul Connect jobspec stanza.
type ConsulConnect struct {
	// Native is true if a service implements Connect directly and does not
	// need a sidecar.
	Native bool

	// SidecarService is non-nil if a service requires a sidecar.
	SidecarService *ConsulSidecarService
}

// Copy the stanza recursively. Returns nil if nil.
func (c *ConsulConnect) Copy() *ConsulConnect {
	if c == nil {
		return nil
	}

	return &ConsulConnect{
		Native:         c.Native,
		SidecarService: c.SidecarService.Copy(),
	}
}

// Equals returns true if the structs are recursively equal.
func (c *ConsulConnect) Equals(o *ConsulConnect) bool {
	if c == nil || o == nil {
		return c == o
	}

	if c.Native != o.Native {
		return false
	}

	return c.SidecarService.Equals(o.SidecarService)
}

// Validate that the Connect stanza has exactly one of Native or sidecar.
func (c *ConsulConnect) Validate() error {
	if c.Native && c.SidecarService != nil {
		return fmt.Errorf("Consul Connect must be native or use a sidecar service; not both")
	}

	if !c.Native && c.SidecarService == nil {
		return fmt.Errorf("Consul Connect must be native or use a sidecar service")
	}

	return nil
}

// ConsulSidecarService represents a Consul Connect SidecarService jobspec
// stanza.
type ConsulSidecarService struct {
	// Port is the service's port that the sidecar will connect to. May be
	// a port label or a literal port number.
	Port string

	// Proxy stanza defining the sidecar proxy configuration.
	Proxy *ConsulProxy
}

// Copy the stanza recursively. Returns nil if nil.
func (s *ConsulSidecarService) Copy() *ConsulSidecarService {
	return &ConsulSidecarService{
		Port:  s.Port,
		Proxy: s.Proxy.Copy(),
	}
}

// Equals returns true if the structs are recursively equal.
func (s *ConsulSidecarService) Equals(o *ConsulSidecarService) bool {
	if s == nil || o == nil {
		return s == o
	}

	if s.Port != o.Port {
		return false
	}

	return s.Proxy.Equals(o.Proxy)
}

// ConsulProxy represents a Consul Connect sidecar proxy jobspec stanza.
type ConsulProxy struct {
	// Upstreams configures the upstream services this service intends to
	// connect to.
	Upstreams []*ConsulUpstream

	// Config is a proxy configuration. It is opaque to Nomad and passed
	// directly to Consul.
	Config map[string]interface{}
}

// Copy the stanza recursively. Returns nil if nil.
func (p *ConsulProxy) Copy() *ConsulProxy {
	if p == nil {
		return nil
	}

	newP := ConsulProxy{}

	if n := len(p.Upstreams); n > 0 {
		newP.Upstreams = make([]*ConsulUpstream, n)

		for i := range p.Upstreams {
			newP.Upstreams[i] = p.Upstreams[i].Copy()
		}
	}

	if n := len(p.Config); n > 0 {
		newP.Config = make(map[string]interface{}, n)

		for k, v := range p.Config {
			newP.Config[k] = v
		}
	}

	return &newP
}

// Equals returns true if the structs are recursively equal.
func (p *ConsulProxy) Equals(o *ConsulProxy) bool {
	if p == nil || o == nil {
		return p == o
	}

	fmt.Println("len-->", p.Upstreams, o.Upstreams)
	if len(p.Upstreams) != len(o.Upstreams) {
		return false
	}

	// Order doesn't matter
OUTER:
	for _, up := range p.Upstreams {
		for _, innerUp := range o.Upstreams {
			if up.Equals(innerUp) {
				// Match; find next upstream
				continue OUTER
			}
		}

		// No match
		return false
	}

	// Avoid nil vs {} differences
	if len(p.Config) != 0 && len(o.Config) != 0 {
		if !reflect.DeepEqual(p.Config, o.Config) {
			return false
		}
	}

	return true
}

// ConsulUpstream represents a Consul Connect upstream jobspec stanza.
type ConsulUpstream struct {
	// DestinationName is the name of the upstream service.
	DestinationName string

	// LocalBindPort is the port the proxy will receive connections for the
	// upstream on.
	LocalBindPort int
}

// Copy the stanza recursively. Returns nil if nil.
func (u *ConsulUpstream) Copy() *ConsulUpstream {
	if u == nil {
		return nil
	}

	return &ConsulUpstream{
		DestinationName: u.DestinationName,
		LocalBindPort:   u.LocalBindPort,
	}
}

// Equals returns true if the structs are recursively equal.
func (u *ConsulUpstream) Equals(o *ConsulUpstream) bool {
	if u == nil || o == nil {
		return u == o
	}

	return (*u) == (*o)
}
