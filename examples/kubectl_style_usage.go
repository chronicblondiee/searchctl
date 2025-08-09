package main

import (
	"fmt"
	"log"

	"github.com/chronicblondiee/searchctl/pkg/client"
)

func main() {
	// Create clientset using kubectl-style architecture
	clientset, err := client.NewClientset()
	if err != nil {
		log.Fatal(err)
	}

	// kubectl-style resource-centric API calls
	fmt.Println("=== Cluster Operations ===")
	health, err := clientset.Cluster().Health()
	if err != nil {
		fmt.Printf("Error getting cluster health: %v\n", err)
	} else {
		fmt.Printf("Cluster Status: %s\n", health.Status)
	}

	info, err := clientset.Cluster().Info()
	if err != nil {
		fmt.Printf("Error getting cluster info: %v\n", err)
	} else {
		fmt.Printf("Cluster Name: %s\n", info.ClusterName)
	}

	fmt.Println("\n=== Index Operations ===")
	indices, err := clientset.Indices().List("logs-*")
	if err != nil {
		fmt.Printf("Error listing indices: %v\n", err)
	} else {
		fmt.Printf("Found %d indices matching 'logs-*'\n", len(indices))
	}

	// Create index with settings
	indexBody := map[string]interface{}{
		"settings": map[string]interface{}{
			"number_of_shards":   1,
			"number_of_replicas": 0,
		},
	}
	err = clientset.Indices().Create("test-index", indexBody)
	if err != nil {
		fmt.Printf("Error creating index: %v\n", err)
	} else {
		fmt.Println("Successfully created test-index")
	}

	fmt.Println("\n=== Index Template Operations ===")
	templates, err := clientset.Indices().Templates().List("logs-*")
	if err != nil {
		fmt.Printf("Error listing templates: %v\n", err)
	} else {
		fmt.Printf("Found %d templates matching 'logs-*'\n", len(templates))
	}

	fmt.Println("\n=== Data Stream Operations ===")
	dataStreams, err := clientset.DataStreams().List("metrics-*")
	if err != nil {
		fmt.Printf("Error listing data streams: %v\n", err)
	} else {
		fmt.Printf("Found %d data streams matching 'metrics-*'\n", len(dataStreams))
	}

	fmt.Println("\n=== Node Operations ===")
	nodes, err := clientset.Nodes().List()
	if err != nil {
		fmt.Printf("Error listing nodes: %v\n", err)
	} else {
		fmt.Printf("Found %d nodes in cluster\n", len(nodes))
	}

	// Backward compatibility - old client interface still works
	fmt.Println("\n=== Backward Compatibility ===")
	oldClient, err := client.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	oldHealth, err := oldClient.ClusterHealth()
	if err != nil {
		fmt.Printf("Error using old client: %v\n", err)
	} else {
		fmt.Printf("Old client works - Cluster Status: %s\n", oldHealth.Status)
	}
}