package helper

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	apiv1 "redis-crd/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
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

func Create(client client.Client, redis *apiv1.Redis, name string) (string, error) {
	if IsExist(name, redis) {
		return "", nil
	}
	pod := &corev1.Pod{}
	pod.Name = redis.Name
	pod.Namespace = redis.Namespace

	pod.Spec.Containers = []corev1.Container{
		{
			Name:            redis.Name,
			Image:           "redis:6.2-alpine3.15",
			ImagePullPolicy: corev1.PullIfNotPresent,
			Ports: []corev1.ContainerPort{
				{
					ContainerPort: int32(redis.Spec.Port),
				},
			},
		},
	}

	return name, client.Create(context.Background(), pod)
}
