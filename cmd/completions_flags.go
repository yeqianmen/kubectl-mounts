package cmd

import (
	"context"
	"os"
	"strings"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// RegisterCompletions all Auto completion flag
func RegisterCompletions(cmd *cobra.Command) {
	cmd.RegisterFlagCompletionFunc("namespace", completeNamespaces)
	cmd.RegisterFlagCompletionFunc("pod", completePods)
	cmd.RegisterFlagCompletionFunc("output", completeOutput)
}

// Auto completion--namespace
func completeNamespaces(_ *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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

	var names []string
	for _, ns := range nsList.Items {
		if strings.HasPrefix(ns.Name, toComplete) {
			names = append(names, ns.Name)
		}
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}

// Auto completion --pod
func completePods(_ *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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

	var names []string
	for _, pod := range pods.Items {
		if strings.HasPrefix(pod.Name, toComplete) {
			names = append(names, pod.Name)
		}
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}

// Auto completion --output
func completeOutput(_ *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	options := []string{"json", "yaml", "table"}
	var filtered []string
	for _, opt := range options {
		if strings.HasPrefix(opt, toComplete) {
			filtered = append(filtered, opt)
		}
	}
	return filtered, cobra.ShellCompDirectiveNoFileComp
}
