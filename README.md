# Virtual Service Route Controller

Construct Istio `VirtualService` resource from some components.

:writing_hand: This is study about custom controller.

## Usage
1. Create `VirtualServiceBase` resource.

```bash
$ cat <<EOF | kubectl apply -f -
apiVersion: virtualservicecomponent.shibataka000.com/v1alpha1
kind: VirtualServiceBase
metadata:
  name: virtualservicebase-sample
spec:
  gateways:
  - gateway
  hosts:
  - '*'
EOF
```

2. Create `HTTPRouteBinding` resources.

```bash
$ cat <<EOF | kubectl apply -f -
apiVersion: virtualservicecomponent.shibataka000.com/v1alpha1
kind: HTTPRouteBinding
metadata:
  name: httproutebinding-sample-1
spec:
  virtualServiceBaseRef:
    apiVersion: virtualservicecomponent.shibataka000.com/v1alpha1
    kind: VirtualServiceBase
    name: virtualservicebase-sample
    namespace: default
  httpRoute:
    match:
    - headers:
        key:
          exact: mysubset-1
    route:
    - destination:
        host: myhost
        subset: mysubset-1
EOF

$ cat <<EOF | kubectl apply -f -
apiVersion: virtualservicecomponent.shibataka000.com/v1alpha1
kind: HTTPRouteBinding
metadata:
  name: httproutebinding-sample-2
spec:
  virtualServiceBaseRef:
    apiVersion: virtualservicecomponent.shibataka000.com/v1alpha1
    kind: VirtualServiceBase
    name: virtualservicebase-sample
    namespace: default
  httpRoute:
    match:
    - headers:
        key:
          exact: mysubset-2
    route:
    - destination:
        host: myhost
        subset: mysubset-2
EOF
```

3. VirtualServiceRouteController construct `VirtualService` resource from `VirtualServiceBase` and `HTTPRouteBinding` . You can see created resource as follows.

```bash
$ kubectl get virtualservices.networking.istio.io virtualservicebase-sample -o yaml
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: virtualservicebase-sample
  namespace: default
spec:
  gateways:
  - gateway
  hosts:
  - '*'
  http:
  - match:
    - headers:
        key:
          exact: mysubset-1
    route:
    - destination:
        host: myhost
        subset: mysubset-1
  - match:
    - headers:
        key:
          exact: mysubset-2
    route:
    - destination:
        host: myhost
        subset: mysubset-2
```

## Requirement
- [Istio](https://istio.io/)

## Install
1. Deploy Istio to your kubernetes cluster. See [Setup](https://istio.io/latest/docs/setup/) more details.
2. Deploy VirtualServiceRouteController and CRDs to your kubernetes cluster.

```bash
TBD
```
