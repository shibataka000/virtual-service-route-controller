package v1alpha1

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
