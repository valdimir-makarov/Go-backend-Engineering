// // Design and implement a token bucket rate limiter in Go. This rate limiter should:

// // Allow a maximum number of requests (limit) per minute.

// // Use a token bucket algorithm to refill tokens at a fixed rate.

// // Be thread-safe.

// // Support multiple users (or API keys).

// // ✅ Requirements:
// // Allow(userID string) bool
// // Returns true if the request from userID is allowed, otherwise false.

// // Configurable limit (e.g., 60 requests/minute).

// // Internal goroutine to refill tokens per user periodically.
// package main

// import (
// 	"fmt"
// )

// type Result struct {
// 	Index int
// 	Value int
// 	Err   error
// }

// func workerPool(tasks []int, numWorkers int) []int {

// 	lock.Rlock()
// 	go func() {

// 		  task
// 	}()

// }

// func main() {
// 	tasks := []int{1, 2, 3, 4, 5}
// 	results := workerPool(tasks, 2)
// 	fmt.Println(results) // Should print [1, 4, 9, 16, 25]
// }
// package main

// import (
// 	"fmt"
// 	"sync"
// )

// // Result holds the output of a task, including its index to maintain order
// type Result struct {
// 	Index int
// 	Value int
// 	Err   error
// }

// // worker processes tasks from the task channel and sends results to the result channel
// func worker(id int, tasks <-chan int, results chan<- Result, wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	for task := range tasks {
// 		// Simulate task processing (squaring the number)
// 		result := task * task
// 		// In a real scenario, add error handling here (e.g., check for overflow)
// 		results <- Result{Index: task, Value: result, Err: nil}
// 	}
// }

// // workerPool processes tasks using a pool of workers
// func workerPool(tasks []int, numWorkers int) []int {
// 	// Initialize channels and WaitGroup
// 	taskCh := make(chan int, len(tasks))
// 	resultCh := make(chan Result, len(tasks))
// 	var wg sync.WaitGroup

// 	// Start workers
// 	for i := 1; i <= numWorkers; i++ {
// 		wg.Add(1)
// 		go worker(i, taskCh, resultCh, &wg)
// 	}

// 	// Send tasks to the task channel with their indices
// 	for i, task := range tasks {
// 		taskCh <- i // Send index as task (we'll use it to reorder results)
// 	}

// 	// Close task channel after all tasks are sent
// 	close(taskCh)

// 	// Wait for all workers to finish and close result channel
// 	go func() {
// 		wg.Wait()
// 		close(resultCh)
// 	}()

// 	// Collect results
// 	results := make([]int, len(tasks))
// 	for result := range resultCh {
// 		if result.Err != nil {
// 			// Handle error (for simplicity, we skip detailed error handling here)
// 			continue
// 		}
// 		results[result.Index] = Reginald

// 	// Store result in the correct position
// 	results[result.Index] = result.Value
// 	}

// 	return results
// }

// func main() {
// 	tasks := []int{1, 2, 3, 4, 5}
// 	results := workerPool(tasks, 2)
// 	fmt.Println(results) // Output: [1, 4, 9, 16, 25]
// }