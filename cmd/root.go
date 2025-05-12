package cmd

import (
	"context"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
	"strings"
)

var (
	namespace  string
	podFilter  string
	kubeconfig string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kubectl-mounts",
	Short: "Show Pod Volumes and VolumeMounts in the cluster",
	Run: func(cmd *cobra.Command, args []string) {
		runMounts()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	rootCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Namespace (default is current namespace)")
	rootCmd.Flags().StringVarP(&podFilter, "pod", "p", "", "Filter by specific Pod name")
	rootCmd.Flags().StringVarP(&kubeconfig, "kubeconfig", "k", "", "Path to kubeconfig file (default $HOME/.kube/config)")
	// Register namespace completion
	rootCmd.RegisterFlagCompletionFunc("namespace", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		config, err := getKubeConfig()
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		nsList, err := clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		names := make([]string, 0, len(nsList.Items))
		for _, ns := range nsList.Items {
			if strings.HasPrefix(ns.Name, toComplete) {
				names = append(names, ns.Name)
			}
		}
		return names, cobra.ShellCompDirectiveNoFileComp
	})
	// Register pod completion
	rootCmd.RegisterFlagCompletionFunc("pod", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		config, err := getKubeConfig()
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		ns := namespace
		if ns == "" {
			nsBytes, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
			if err != nil {
				ns = "default"
			} else {
				ns = strings.TrimSpace(string(nsBytes))
			}
		}

		pods, err := clientset.CoreV1().Pods(ns).List(context.Background(), metav1.ListOptions{})
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		names := make([]string, 0, len(pods.Items))
		for _, pod := range pods.Items {
			if strings.HasPrefix(pod.Name, toComplete) {
				names = append(names, pod.Name)
			}
		}
		return names, cobra.ShellCompDirectiveNoFileComp
	})
}

func runMounts() {
	config, err := getKubeConfig()
	if err != nil {
		fmt.Println("Failed to get Kubernetes config:", err)
		os.Exit(1)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println("Failed to create Kubernetes client:", err)
		os.Exit(1)
	}

	ns := namespace
	if ns == "" {
		nsBytes, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
		if err != nil {
			ns = "default"
		} else {
			ns = strings.TrimSpace(string(nsBytes))
		}
	}

	pods, err := clientset.CoreV1().Pods(ns).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println("Failed to list Pods:", err)
		os.Exit(1)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Pod", "Container", "Volume Name", "MountPath", "Volume Type"})
	table.SetAutoFormatHeaders(false)
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)

	for _, pod := range pods.Items {
		if podFilter != "" && pod.Name != podFilter {
			continue
		}

		// Build Volume name and type mapping
		volumeTypeMap := make(map[string]string)
		for _, v := range pod.Spec.Volumes {
			volumeTypeMap[v.Name] = describeVolumeSource(v)
		}

		for _, c := range pod.Spec.Containers {
			for _, m := range c.VolumeMounts {
				volumeType := volumeTypeMap[m.Name]
				table.Append([]string{
					pod.Name,
					c.Name,
					m.Name,
					m.MountPath,
					volumeType,
				})
			}
		}
	}

	table.Render()
}

func describeVolumeSource(v corev1.Volume) string {
	switch {
	case v.EmptyDir != nil:
		return "EmptyDir"
	case v.HostPath != nil:
		return "HostPath"
	case v.PersistentVolumeClaim != nil:
		return "PVC"
	case v.ConfigMap != nil:
		return "ConfigMap"
	case v.Secret != nil:
		return "Secret"
	case v.Projected != nil:
		return "Projected"
	default:
		return "Other"
	}
}

func getKubeConfig() (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	defaultKubeConfig := filepath.Join(home, ".kube", "config")
	return clientcmd.BuildConfigFromFlags("", defaultKubeConfig)
}
