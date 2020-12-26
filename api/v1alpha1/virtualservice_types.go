/*
This file is based on https://github.com/istio/api/blob/master/networking/v1beta1/virtual_service.pb.go

I'd like to refer virtual_service.pb.go directly, but controller-gen can't treat it because of lack of json tag.
So I copy and modify some structures which are needed.
*/

package v1alpha1

// VirtualService is configuration affecting traffic routing.
type VirtualService struct {
	Hosts    []string `json:"hosts,omitempty"`
	Gateways []string `json:"gateways,omitempty"`
}

// HTTPRoute describes match conditions and actions for routing HTTP/1.1, HTTP2, and gRPC traffic.
type HTTPRoute struct {
	Match []*HTTPMatchRequest     `json:"match,omitempty"`
	Route []*HTTPRouteDestination `json:"route,omitempty"`
}

// HTTPMatchRequest specifies a set of criterion to be met in order for the rule to be applied to the HTTP request.
type HTTPMatchRequest struct {
	Headers map[string]*StringMatch `json:"headers,omitempty"`
}

// StringMatch describes how to match a given string in HTTP headers. Match is case-sensitive.
type StringMatch struct {
	Exact string `json:"exact,omitempty"`
}

// HTTPRouteDestination has routing rules which are associated with one or more service versions.
type HTTPRouteDestination struct {
	Destination *Destination `json:"destination,omitempty"`
}

// Destination indicates the network addressable service to which the request/connection will be sent after processing a routing rule.
type Destination struct {
	Host   string `json:"host,omitempty"`
	Subset string `json:"subset,omitempty"`
}
