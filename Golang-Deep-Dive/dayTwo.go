// package main

// import (
// 	"fmt"
// 	"sort"
// )

// // You are building a simple task scheduler for a small business application. The scheduler accepts tasks with a name and a priority (an integer from 1 to 10, where 1 is the highest priority). Your program should allow adding tasks to a queue and retrieving the highest-priority task to execute. If multiple tasks have the same priority, the oldest task (first added) should be executed first.

// // Requirements:

// // Implement a Task struct with fields for the task name and priority.
// // Create a TaskScheduler struct that can:
// // Add a task to the queue.
// // Retrieve and remove the highest-priority task (lowest priority number).
// // Return the current list of tasks in the queue, sorted by priority and insertion order.
// // Write a simple main function to demonstrate adding tasks and retrieving the highest-priority task.
// // Constraints:

// // Do not use external libraries beyond Go's standard library.
// // The solution should handle cases where the queue is empty.
// // Keep the code clean, readable, and well-commented.

// type Task struct {
// 	name string
// 	pr   int64
// }

// var s1 = make([]*Task, 0, 100)

// func (T *Task) TaskScheduler(name string, pr int64) string {

// 	// enque a Task

// 	s1 = append(s1, &Task{
// 		name: name,
// 		pr:   pr,
// 	})
// 	if s1 == nil {
// 		return "failed to inset"
// 	}
// 	return "inserted successfully"
// }

// func (T *Task) DeleteTask(selectedTask string) string {

// 	for index, value := range s1 {
// 		if value.name == selectedTask {
// 			returned := value
// 			s1 = append(s1[:index], s1[index+1:]...)

// 			return returned.name
// 		}
// 	}
// 	return "could fint the Task"
// }
// func (T *Task) Get() *Task {
// 	sort.Slice(s1, func(i, j int) bool {
// 		return s1[i].pr < s1[j].pr
// 	})

// 	return s1[len(s1)-1]
// }
// func main() {
// 	name1 := "Task 1"
// 	name2 := "Task 2"
// 	name3 := "Task 3"
// 	name4 := "Task 4"
// 	name5 := "Task 5"
// 	name6 := "Task 6"
// 	name7 := "Task 7"
// 	name8 := "Task 8"
// 	name9 := "Task 9"
// 	name10 := "Task 10"

// 	t := &Task{}

// 	// Add 10 tasks manually
// 	result1 := t.TaskScheduler(name1, 1)
// 	result2 := t.TaskScheduler(name2, 2)
// 	result3 := t.TaskScheduler(name3, 3)
// 	result4 := t.TaskScheduler(name4, 4)
// 	result5 := t.TaskScheduler(name5, 5)
// 	result6 := t.TaskScheduler(name6, 6)
// 	result7 := t.TaskScheduler(name7, 7)
// 	result8 := t.TaskScheduler(name8, 8)
// 	result9 := t.TaskScheduler(name9, 9)
// 	result10 := t.TaskScheduler(name10, 1)

// 	// Print all results
// 	fmt.Println("hey", result1)
// 	fmt.Println("hey", result2)
// 	fmt.Println("hey", result3)
// 	fmt.Println("hey", result4)
// 	fmt.Println("hey", result5)
// 	fmt.Println("hey", result6)
// 	fmt.Println("hey", result7)
// 	fmt.Println("hey", result8)
// 	fmt.Println("hey", result9)
// 	fmt.Println("hey", result10)

// 	// Optional: print all tasks in queue
// 	fmt.Println("\nTasks in queue:")
// 	for i, task := range s1 {
// 		fmt.Printf("%d: Name=%s, Priority=%d\n", i+1, task.name, task.pr)
// 	}
// 	// deleted := t.DeleteTask("bubun")
// 	// fmt.Printf("%v", deleted)
// 	get := t.Get()
// 	fmt.Printf("%v", get)
// }
