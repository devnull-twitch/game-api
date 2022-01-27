package k8s

import (
	"context"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetPortFinder(clientset *kubernetes.Clientset) func(string) int {
	return func(zone string) int {
		k8sZone := strings.ReplaceAll(zone, "_", "-")

		services, err := clientset.CoreV1().Services("default").List(context.Background(), metav1.ListOptions{
			LabelSelector: fmt.Sprintf("gamezone=%s", k8sZone),
		})
		if err != nil {
			logrus.WithError(err).Error("unable to load services")
			return 0
		}

		if len(services.Items) <= 0 {
			logrus.WithError(err).Warn("no matching service")
			return 0
		}

		return int(services.Items[0].Spec.Ports[0].NodePort)
	}
}
