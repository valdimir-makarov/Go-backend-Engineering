package main

import (
	"fmt"
	"log"
	"os/exec"
)

type ContainerPool struct {
	containers []string
	available  chan string
}

func NewContainerPool(containers []string) *ContainerPool {
	pool := &ContainerPool{
		containers: containers,
		available:  make(chan string, len(containers)),
	}

	for _, c := range containers {
		pool.available <- c
	}

	return pool
}

// Create & start N containers
func initContainers(N int) []string {
	var containers []string
	for i := 1; i <= N; i++ {
		containerName := fmt.Sprintf("code-runner-%d", i)

		// Check if container already exists
		cmd := exec.Command("docker", "inspect", "--type=container", containerName)
		if err := cmd.Run(); err == nil {
			log.Printf("Container %s already exists.", containerName)
		} else {
			// Create and start the container
			cmd := exec.Command("docker", "run", "-d", "--name", containerName, "python:3.11", "tail", "-f", "/dev/null")
			if output, err := cmd.CombinedOutput(); err != nil {
				log.Fatalf("Failed to start container %s: %v\n%s", containerName, err, string(output))
			} else {
				log.Printf("Started container: %s", containerName)
			}
		}

		containers = append(containers, containerName)
	}
	return containers
}

var pool *ContainerPool

func main() {
	// Initialize containers at startup
	containerNames := initContainers(5) // for example, 5 containers
	pool = NewContainerPool(containerNames)

	// Now you can start your HTTP server here, using the pool
	// e.g., http.ListenAndServe(":8080", yourHandler)
}
