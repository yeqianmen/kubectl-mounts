package main

import (
	"context"
	"fmt"
	"github.com/olekukonko/tablewriter"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

func main() {
	// 加载 kubeconfig
	home := os.Getenv("HOME")
	kubeconfig := filepath.Join(home, ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	namespace := getCurrentNamespace()

	pods, err := clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	if len(pods.Items) == 0 {
		fmt.Printf("No pods found in namespace: %s\n", namespace)
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Pod", "Container", "Volume Name", "MountPath", "Volume Type"})

	for _, pod := range pods.Items {
		volTypeMap := make(map[string]string)
		for _, vol := range pod.Spec.Volumes {
			typ := getVolumeType(vol)
			volTypeMap[vol.Name] = typ
		}

		for _, container := range pod.Spec.Containers {
			for _, mount := range container.VolumeMounts {
				row := []string{
					pod.Name,
					container.Name,
					mount.Name,
					mount.MountPath,
					volTypeMap[mount.Name],
				}
				table.Append(row)
			}
		}
	}

	table.Render()
}

func getCurrentNamespace() string {
	namespace := "default"
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	config, err := rules.Load()
	if err == nil {
		currentContext := config.CurrentContext
		if ctx, ok := config.Contexts[currentContext]; ok {
			if ctx.Namespace != "" {
				namespace = ctx.Namespace
			}
		}
	}
	return namespace
}

func getVolumeType(vol v1.Volume) string {
	switch {
	case vol.EmptyDir != nil:
		return "EmptyDir"
	case vol.ConfigMap != nil:
		return fmt.Sprintf("ConfigMap(%s)", vol.ConfigMap.Name)
	case vol.Secret != nil:
		return fmt.Sprintf("Secret(%s)", vol.Secret.SecretName)
	case vol.PersistentVolumeClaim != nil:
		return fmt.Sprintf("PVC(%s)", vol.PersistentVolumeClaim.ClaimName)
	default:
		return "Other"
	}
}
