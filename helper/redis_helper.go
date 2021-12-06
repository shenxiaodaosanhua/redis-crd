package helper

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	apiv1 "redis-crd/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func Create(client client.Client, redis *apiv1.Redis) error {
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

	return client.Create(context.Background(), pod)
}
