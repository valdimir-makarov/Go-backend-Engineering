package pool

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
		available:  make(chan string, 10), // Increased capacity
	}

	for _, c := range containers {
		pool.available <- c
		log.Printf("Added to pool: %s", c)
	}
	return pool
}

// Clean up existing containers
func CleanupContainers(names []string) {
	for _, name := range names {
		cmd := exec.Command("docker", "inspect", "--type=container", name)
		if err := cmd.Run(); err == nil {
			cmd := exec.Command("docker", "rm", "-f", name)
			if output, err := cmd.CombinedOutput(); err != nil {
				log.Printf("Failed to remove container %s: %v\n%s", name, err, string(output))
			} else {
				log.Printf("Removed container: %s", name)
			}
		}
	}
}

// Create & start N containers
func InitContainers(N int) []string {
	var containers []string
	for i := 1; i <= N; i++ {
		containerName := fmt.Sprintf("code-runner-%d", i)
		cmd := exec.Command("docker", "inspect", "--type=container", containerName)
		if err := cmd.Run(); err == nil {
			log.Printf("Container %s already exists.", containerName)
		} else {
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

func (p *ContainerPool) GetContainer() string {
	c := <-p.available
	fmt.Println("Got container:", c, "Channel len:", len(p.available))
	return c
}

func (p *ContainerPool) ReleaseContainer(name string) {
	select {
	case p.available <- name:
		fmt.Println("Releasing container:", name, "Channel len:", len(p.available))
	default:
		fmt.Println("Channel full, dropping:", name)
	}
}

// func main() {
// 	start := time.Now()

// 	var containerNames []string
// 	for i := 1; i <= 5; i++ {
// 		containerNames = append(containerNames, fmt.Sprintf("code-runner-%d", i))
// 	}
// 	CleanupContainers(containerNames)

// 	container := InitContainers(5)
// 	fmt.Println("Initialized containers:", container)
// 	if len(container) == 0 {
// 		log.Fatal("No containers initialized.")
// 	}
// 	pool := NewContainerPool(container)
// 	for i := 0; i < 5; i++ {
// 		c := pool.GetContainer()
// 		fmt.Printf("Using container: %s\n", c)
// 		pool.ReleaseContainer(c)
// 	}

// 	duration := time.Since(start)
// 	fmt.Printf("Total execution time: %s\n", duration)
// }
