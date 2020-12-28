/*
This file is based on https://github.com/istio/api/blob/master/networking/v1beta1/virtual_service.pb.go

I'd like to refer virtual_service.pb.go directly, but controller-gen can't treat it because of lack of json tag.
So I copy and modify some structures which are needed.
*/

package v1alpha1

import (
	networkingv1beta1 "istio.io/api/networking/v1beta1"
)

// VirtualService is configuration affecting traffic routing.
type VirtualService struct {
	Hosts    []string `json:"hosts,omitempty"`
	Gateways []string `json:"gateways,omitempty"`
}

// IstioAPI convert to VirtualService struct defined in istio.io/api
func (a *VirtualService) IstioAPI() *networkingv1beta1.VirtualService {
	return &networkingv1beta1.VirtualService{
		Hosts:    a.Hosts,
		Gateways: a.Gateways,
	}
}

// HTTPRoute describes match conditions and actions for routing HTTP/1.1, HTTP2, and gRPC traffic.
type HTTPRoute struct {
	Match []*HTTPMatchRequest     `json:"match,omitempty"`
	Route []*HTTPRouteDestination `json:"route,omitempty"`
}

// IstioAPI convert to HTTPRoute struct defined in istio.io/api
func (a *HTTPRoute) IstioAPI() *networkingv1beta1.HTTPRoute {
	matches := []*networkingv1beta1.HTTPMatchRequest{}
	for _, match := range a.Match {
		matches = append(matches, match.IstioAPI())
	}
	routes := []*networkingv1beta1.HTTPRouteDestination{}
	for _, route := range a.Route {
		routes = append(routes, route.IstioAPI())
	}
	return &networkingv1beta1.HTTPRoute{
		Match: matches,
		Route: routes,
	}
}

// HTTPMatchRequest specifies a set of criterion to be met in order for the rule to be applied to the HTTP request.
type HTTPMatchRequest struct {
	Headers map[string]*StringMatch `json:"headers,omitempty"`
}

// IstioAPI convert to HTTPMatchRequest struct defined in istio.io/api
func (a *HTTPMatchRequest) IstioAPI() *networkingv1beta1.HTTPMatchRequest {
	headers := map[string]*networkingv1beta1.StringMatch{}
	for key, value := range a.Headers {
		headers[key] = value.IstioAPI()
	}
	return &networkingv1beta1.HTTPMatchRequest{
		Headers: headers,
	}
}

// StringMatch describes how to match a given string in HTTP headers. Match is case-sensitive.
type StringMatch struct {
	Exact string `json:"exact,omitempty"`
}

// IstioAPI convert to StringMatch struct defined in istio.io/api
func (a *StringMatch) IstioAPI() *networkingv1beta1.StringMatch {
	return &networkingv1beta1.StringMatch{
		MatchType: &networkingv1beta1.StringMatch_Exact{
			Exact: a.Exact,
		},
	}
}

// HTTPRouteDestination has routing rules which are associated with one or more service versions.
type HTTPRouteDestination struct {
	Destination *Destination `json:"destination,omitempty"`
}

// IstioAPI convert to HTTPRouteDestination struct defined in istio.io/api
func (a *HTTPRouteDestination) IstioAPI() *networkingv1beta1.HTTPRouteDestination {
	return &networkingv1beta1.HTTPRouteDestination{
		Destination: a.Destination.IstioAPI(),
	}
}

// Destination indicates the network addressable service to which the request/connection will be sent after processing a routing rule.
type Destination struct {
	Host   string `json:"host,omitempty"`
	Subset string `json:"subset,omitempty"`
}

// IstioAPI convert to Destination struct defined in istio.io/api
func (a *Destination) IstioAPI() *networkingv1beta1.Destination {
	return &networkingv1beta1.Destination{
		Host:   a.Host,
		Subset: a.Subset,
	}
}
