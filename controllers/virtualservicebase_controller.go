/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"reflect"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	virtualservicecomponentv1alpha1 "github.com/shibataka000/virtual-service-route-controller/api/v1alpha1"
	networkingv1beta1 "istio.io/api/networking/v1beta1"
	networkingv1beta1client "istio.io/client-go/pkg/apis/networking/v1beta1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	virtualServiceOwnerKey = ".metadata.controller"
	apiGVStr               = virtualservicecomponentv1alpha1.GroupVersion.String()
)

// VirtualServiceBaseReconciler reconciles a VirtualServiceBase object
type VirtualServiceBaseReconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=virtualservicecomponent.shibataka000.com,resources=virtualservicebases,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=virtualservicecomponent.shibataka000.com,resources=virtualservicebases/status,verbs=get;update;patch

func (r *VirtualServiceBaseReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("virtualservicebase", req.NamespacedName)

	var virtualServiceBase virtualservicecomponentv1alpha1.VirtualServiceBase
	if err := r.Get(ctx, req.NamespacedName, &virtualServiceBase); err != nil {
		log.Error(err, "unable to fetch VirtualServiceBase")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if err := r.cleanupOwnedResources(ctx, log, &virtualServiceBase); err != nil {
		log.Error(err, "failed to clean up old VirtualService resources from this VirtualServiceBase")
		return ctrl.Result{}, err
	}

	virtualService := &networkingv1beta1client.VirtualService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      virtualServiceBase.Name,
			Namespace: virtualServiceBase.Namespace,
		},
	}

	if _, err := ctrl.CreateOrUpdate(ctx, r.Client, virtualService, func() error {
		v := virtualServiceBase.Spec.DeepCopy().IstioAPI()

		// hosts
		if v.Hosts == nil {
			err := fmt.Errorf("virtual service must have at least one host")
			log.Error(err, "virtual service must have at least one host")
			return err
		}
		if !reflect.DeepEqual(virtualService.Spec.Hosts, v.Hosts) {
			virtualService.Spec.Hosts = v.Hosts
		}

		// gateways
		if !reflect.DeepEqual(virtualService.Spec.Gateways, v.Gateways) {
			virtualService.Spec.Gateways = v.Gateways
		}

		// http routes
		var httpRouteBindings virtualservicecomponentv1alpha1.HTTPRouteBindingList
		if err := r.List(ctx, &httpRouteBindings, client.InNamespace(virtualServiceBase.Namespace)); err != nil {
			log.Error(err, "unable to fetch HTTPRouteBindings")
			return err
		}
		httpRoutes := []*networkingv1beta1.HTTPRoute{}
		for _, httpRouteBinding := range httpRouteBindings.Items {
			if httpRouteBinding.Spec.VirtualServiceBaseRef.IsReference(&virtualServiceBase) {
				httpRoutes = append(httpRoutes, httpRouteBinding.Spec.HTTPRoute.DeepCopy().IstioAPI())
			}
		}
		if len(httpRoutes) == 0 {
			err := fmt.Errorf("http, tcp or tls must be provided in virtual service")
			log.Error(err, "http, tcp or tls must be provided in virtual service")
			return err
		}
		if !reflect.DeepEqual(virtualService.Spec.Http, httpRoutes) {
			virtualService.Spec.Http = httpRoutes
		}

		if err := ctrl.SetControllerReference(&virtualServiceBase, virtualService, r.Scheme); err != nil {
			log.Error(err, "unable to set ownerReference from VirtualServiceBase to VirtualService")
			return err
		}

		log.Info(fmt.Sprintf("creating or updating VirtualService resource: %s/%s", virtualService.Namespace, virtualService.Name))

		return nil
	}); err != nil {
		log.Error(err, "failed to create or update VirtualService from this VirtualServiceBase")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *VirtualServiceBaseReconciler) SetupWithManager(mgr ctrl.Manager) error {
	log := r.Log
	ctx := context.Background()

	if err := mgr.GetFieldIndexer().IndexField(ctx, &networkingv1beta1client.VirtualService{}, virtualServiceOwnerKey, func(rawObj client.Object) []string {
		virtualService, ok := rawObj.(*networkingv1beta1client.VirtualService)
		if !ok {
			err := fmt.Errorf("unable to type assertion from rawObj to VirtualService")
			log.Error(err, "unable to type assertion from rawObj to VirtualService")
			return nil
		}
		owner := metav1.GetControllerOf(virtualService)
		if owner == nil {
			return nil
		}
		if owner.APIVersion != apiGVStr || owner.Kind != "VirtualService" {
			return nil
		}
		return []string{owner.Name}
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&virtualservicecomponentv1alpha1.VirtualServiceBase{}).
		Owns(&networkingv1beta1client.VirtualService{}).
		Watches(
			&source.Kind{
				Type: &virtualservicecomponentv1alpha1.HTTPRouteBinding{},
			},
			handler.EnqueueRequestsFromMapFunc(func(rawObj client.Object) []reconcile.Request {
				httpRouteBinding, ok := rawObj.(*virtualservicecomponentv1alpha1.HTTPRouteBinding)
				if !ok {
					err := fmt.Errorf("unable to type assertion from rawObj to HTTPRouteBinding")
					log.Error(err, "unable to type assertion from rawObj to HTTPRouteBinding")
					return nil
				}
				return []reconcile.Request{
					{
						NamespacedName: types.NamespacedName{
							Name:      httpRouteBinding.Spec.VirtualServiceBaseRef.Name,
							Namespace: httpRouteBinding.Spec.VirtualServiceBaseRef.Namespace,
						},
					},
				}
			}),
		).
		Complete(r)
}

func (r *VirtualServiceBaseReconciler) cleanupOwnedResources(ctx context.Context, log logr.Logger, virtualServiceBase *virtualservicecomponentv1alpha1.VirtualServiceBase) error {
	var virtualServices networkingv1beta1client.VirtualServiceList
	if err := r.List(ctx, &virtualServices, client.InNamespace(virtualServiceBase.Namespace), client.MatchingFields{virtualServiceOwnerKey: virtualServiceBase.Name}); err != nil {
		return err
	}

	for _, virtualService := range virtualServices.Items {
		if virtualService.Name == virtualServiceBase.Name {
			continue
		}
		if err := r.Delete(ctx, &virtualService); err != nil {
			log.Error(err, "failed to delete VirtualService resource")
			return err
		}
		log.Info(fmt.Sprintf("delete VirtualService resource: %s/%s", virtualService.Namespace, virtualService.Name))
		r.Recorder.Eventf(virtualServiceBase, corev1.EventTypeNormal, "Deleted", "Delete VirtualService %s/%s", virtualService.Namespace, virtualService.Name)
	}

	return nil
}
