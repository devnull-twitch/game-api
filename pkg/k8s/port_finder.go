package k8s

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetPortFinder(clientset *kubernetes.Clientset) func(string) int {
	return func(zone string) int {
		pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{
			LabelSelector: "app=gameserver",
		})
		if err != nil {
			return 0
		}

		for _, p := range pods.Items {
			if p.GetLabels()["app"] == "gameserver" && p.GetLabels()["gamezone"] == zone {
				services, err := clientset.CoreV1().Services("").List(context.Background(), metav1.ListOptions{
					LabelSelector: fmt.Sprintf("gamezone=%s", zone),
				})
				if err != nil {
					return 0
				}
				return int(services.Items[0].Spec.Ports[0].NodePort)
			}
		}

		return 0
	}
}
