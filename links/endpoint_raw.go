package links

import (
	"fmt"
	"net"

	"github.com/srl-labs/containerlab/utils"
)

// EndpointRaw is the raw (string) representation of an endpoint as defined in the topology file
// for a given link definition.
type EndpointRaw struct {
	Node  string `yaml:"node"`
	Iface string `yaml:"interface"`
	MAC   string `yaml:"mac,omitempty"`
}

// NewEndpointRaw creates a new EndpointRaw struct.
func NewEndpointRaw(node, nodeIf, Mac string) *EndpointRaw {
	return &EndpointRaw{
		Node:  node,
		Iface: nodeIf,
		MAC:   Mac,
	}
}

// Resolve resolves the EndpointRaw into an Endpoint interface that is implemented
// by a concrete endpoint struct such as EndpointBridge, EndpointHost, EndpointVeth.
// The type of an endpoint is determined by the node it belongs to.
// Resolving a raw endpoint adds an associated Link and Node to the endpoint.
// It also adds the endpoint to the node.
func (er *EndpointRaw) Resolve(params *ResolveParams, l Link) (Endpoint, error) {
	// check if the referenced node does exist
	node, exists := params.Nodes[er.Node]
	if !exists {
		return nil, fmt.Errorf("unable to find node %s", er.Node)
	}

	genericEndpoint := NewEndpointGeneric(node, er.Iface)
	genericEndpoint.Link = l

	var err error
	if er.MAC == "" {
		// if mac is not present generate one
		genericEndpoint.MAC, err = utils.GenMac(ClabOUI)
		if err != nil {
			return nil, err
		}
	} else {
		// if MAC is present, set it
		m, err := net.ParseMAC(er.MAC)
		if err != nil {
			return nil, err
		}
		genericEndpoint.MAC = m
	}

	var e Endpoint

	switch node.GetLinkEndpointType() {
	case LinkEndpointTypeBridge:
		e = &EndpointBridge{
			EndpointGeneric: *genericEndpoint,
		}
	case LinkEndpointTypeHost:
		e = &EndpointHost{
			EndpointGeneric: *genericEndpoint,
		}
	case LinkEndpointTypeVeth:
		e = &EndpointVeth{
			EndpointGeneric: *genericEndpoint,
		}
	}

	// also add the endpoint to the node
	node.AddEndpoint(e)

	return e, nil
}
