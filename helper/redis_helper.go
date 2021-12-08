package helper

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	apiv1 "redis-crd/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func GetRedisPodNames(redis *apiv1.Redis) []string {
	names := make([]string, redis.Spec.Num)
	for i := 0; i < redis.Spec.Num; i++ {
		names[i] = fmt.Sprintf("%s-%d", redis.Name, i)
	}
	return names
}

func IsExist(name string, redis *apiv1.Redis) bool {
	for _, p := range redis.Finalizers {
		if name == p {
			return true
		}
	}

	return false
}

func IsExistPod(podName string, redis *apiv1.Redis, client client.Client) bool {
	err := client.Get(context.TODO(), types.NamespacedName{
		Namespace: redis.Namespace,
		Name:      podName,
	}, &corev1.Pod{})
	if err != nil {
		return false
	}

	return true
}

func Create(client client.Client, redis *apiv1.Redis, name string, schema *runtime.Scheme) (string, error) {
	if IsExistPod(name, redis, client) {
		return "", nil
	}
	pod := &corev1.Pod{}
	pod.Name = name
	pod.Namespace = redis.Namespace

	pod.Spec.Containers = []corev1.Container{
		{
			Name:            name,
			Image:           "redis:6.2-alpine3.15",
			ImagePullPolicy: corev1.PullIfNotPresent,
			Ports: []corev1.ContainerPort{
				{
					ContainerPort: int32(redis.Spec.Port),
				},
			},
		},
	}

	err := controllerutil.SetControllerReference(redis, pod, schema)
	if err != nil {
		return "", err
	}

	err = client.Create(context.TODO(), pod)
	if err != nil {
		return "", err
	}
	fmt.Println("创建对象：", name)
	return name, nil
}
