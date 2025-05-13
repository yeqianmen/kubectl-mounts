package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
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
	output     string
)

type MountInfo struct {
	PodName    string `yaml:"pod"`
	Container  string `yaml:"container"`
	VolumeName string `yaml:"volume"`
	MountPath  string `yaml:"mountPath"`
	VolumeType string `yaml:"volumeType"`
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kubectl-mounts",
	Short: "Show Pod Volumes and VolumeMounts in the cluster",
	Run: func(cmd *cobra.Command, args []string) {
		runMounts(cmd)
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
	rootCmd.Flags().StringVarP(&output, "output", "o", "", "Output format: table|yaml|json(default table)")
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

func runMounts(cmd *cobra.Command) {
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
	var results []MountInfo

	for _, pod := range pods.Items {
		if podFilter != "" && pod.Name != podFilter {
			continue
		}

		volumeTypeMap := make(map[string]string)
		for _, v := range pod.Spec.Volumes {
			volumeTypeMap[v.Name] = describeVolumeSource(v)
		}

		for _, c := range pod.Spec.Containers {
			for _, m := range c.VolumeMounts {
				results = append(results, MountInfo{
					PodName:    pod.Name,
					Container:  c.Name,
					VolumeName: m.Name,
					MountPath:  m.MountPath,
					VolumeType: volumeTypeMap[m.Name],
				})
			}
		}
	}

	// Determine output format: YAML, JSON, or table (default).
	outputFormat, _ := cmd.Flags().GetString("output")
	switch outputFormat {
	case "yaml":
		out, err := yaml.Marshal(results)
		if err != nil {
			fmt.Println("Failed to marshal YAML:", err)
			os.Exit(1)
		}
		fmt.Println(string(out))
	case "json":
		out, err := json.MarshalIndent(results, "", "  ")
		if err != nil {
			fmt.Println("Failed to marshal JSON:", err)
			os.Exit(1)
		}
		fmt.Println(string(out))
	default:
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Pod", "Container", "Volume Name", "MountPath", "Volume Type"})
		table.SetAutoFormatHeaders(false)
		table.SetAutoMergeCells(true)
		table.SetRowLine(true)

		for _, item := range results {
			table.Append([]string{
				item.PodName,
				item.Container,
				item.VolumeName,
				item.MountPath,
				item.VolumeType,
			})
		}
		table.Render()
	}
}

func describeVolumeSource(v corev1.Volume) string {
	switch {
	case v.EmptyDir != nil:
		return "EmptyDir"
	case v.HostPath != nil:
		return "HostPath"
	case v.PersistentVolumeClaim != nil:
		return fmt.Sprintf("PVC(%s)", v.PersistentVolumeClaim.ClaimName)
	case v.ConfigMap != nil:
		return fmt.Sprintf("ConfigMap(%s)", v.ConfigMap.Name)
	case v.Secret != nil:
		return fmt.Sprintf("Secret(%s)", v.Secret.SecretName)
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
