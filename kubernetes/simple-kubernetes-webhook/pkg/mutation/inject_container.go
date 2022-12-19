package mutation

import (
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

type injectContainer struct {
	Logger logrus.FieldLogger
}

var _ podMutator = (*injectContainer)(nil)

func (ic injectContainer) Name() string {
	return "inject_container"
}

func (ic injectContainer) Mutate(pod *corev1.Pod) (*corev1.Pod, error) {
	ic.Logger = ic.Logger.WithField("mutation", ic.Name())
	mpod := pod.DeepCopy()

	// build out container slice
	containers := []corev1.Container{{
		Name:  "onecer-container",
		Image: "nginx:latest",
	}}

	// inject containers into pod
	for _, container := range containers {
		ic.Logger.Debugf("pod container injected %s", container)
		injectPodContainer(mpod, container)
	}

	return mpod, nil
}

// injectPodContainer injects a container in a pod
func injectPodContainer(pod *corev1.Pod, container corev1.Container) {
	pod.Spec.Containers = append(pod.Spec.Containers, container)
}

// injectInitContainer inject init container in a pod
func injectInitContainer(pod *corev1.Pod, container corev1.Container) {
	pod.Spec.InitContainers = append(pod.Spec.InitContainers, container)
}
