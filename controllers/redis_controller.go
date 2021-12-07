/*
Copyright 2021.

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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"redis-crd/helper"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	myappv1 "redis-crd/api/v1"
)

// RedisReconciler reconciles a Redis object
type RedisReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=myapp.ipicture.vip,resources=redis,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=myapp.ipicture.vip,resources=redis/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=myapp.ipicture.vip,resources=redis/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Redis object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *RedisReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	redis := &myappv1.Redis{}
	if err := r.Get(ctx, req.NamespacedName, redis); err != nil {
		return ctrl.Result{}, err
	} else {
		if !redis.DeletionTimestamp.IsZero() {
			return ctrl.Result{}, r.ClearRedis(ctx, redis)
		}
		names := helper.GetRedisPodNames(redis)
		isAppend := false
		for _, name := range names {
			n, err := helper.Create(r.Client, redis, name)
			if err != nil {
				return ctrl.Result{}, err
			}
			if n == "" {
				continue
			}
			redis.Finalizers = append(redis.Finalizers, n)
			isAppend = true
		}
		if isAppend {
			err = r.Client.Update(ctx, redis)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	return ctrl.Result{}, nil
}

func (r *RedisReconciler) ClearRedis(ctx context.Context, redis *myappv1.Redis) error {
	list := redis.Finalizers
	for _, name := range list {
		err := r.Client.Delete(ctx, &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: redis.Namespace},
		})
		if err != nil {
			return err
		}
	}

	redis.Finalizers = []string{}
	return r.Client.Update(ctx, redis)
}

// SetupWithManager sets up the controller with the Manager.
func (r *RedisReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&myappv1.Redis{}).
		Complete(r)
}
