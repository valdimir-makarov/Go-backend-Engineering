package main

// Awesome, I love the challenge of ramping up the difficulty! Since you’ve completed the task scheduler and I’ve given you an inventory management problem for tomorrow, I’ll design a new problem for the day after (Day 3) that’s slightly harder, focusing on a real-world scenario with increased complexity while still being achievable for a junior Go developer in under 30 minutes. The problem will test additional Go skills, introduce a bit more logic, and build on the concepts from the previous problems. I’ll also outline a plan to make future problems progressively harder.

// Day 3 Problem: Event Logger with Persistence
// You’re building an event logging system for a small application that tracks user actions (e.g., login, logout, purchase). Each event has a timestamp, action type, and user ID. The system should allow logging events, retrieving events for a specific user, and saving/loading the event log to/from a file for persistence. The challenge includes handling file I/O and ensuring thread-safe access to the log (introducing basic concurrency).

// type Event struct {
// 	Timestamp time.Time
// 	Action    string
// 	User      int64
// }

// var slice = make([]*Event, 0)
// var mu sync.Mutex

// func CreateEvents(action string, user int64) error {
// 	mu.Lock()
// 	defer mu.Unlock()
// 	slice = append(slice, &Event{
// 		Timestamp: time.Now(), // Set timestamp to current time
// 		Action:    action,
// 		User:      user,
// 	})
// 	data, err := json.MarshalIndent(slice, "", "  ")

// 	if err != nil {
// 		panic(err)
// 	}
// 	err = os.WriteFile("output12.json", data, 0644)
// 	if err != nil {
// 		fmt.Printf(err.Error())
// 	}
// 	return nil
// }
// func main() {
// 	CreateEvents("bubun", 49)
// 	for i, e := range slice {
// 		fmt.Printf("Event %d: %+v\n", i, e)
// 	}
// }

// Requirements:

// Define an Event struct with fields for timestamp (time.Time), action (string), and user ID (string).
// Create an EventLogger struct that can:
// Log a new event with the current timestamp.
// Retrieve all events for a given user ID, sorted by timestamp (oldest first).
// Save the entire event log to a JSON file.
// Load the event log from a JSON file, appending to existing events.
// Ensure thread-safe access to the event log using a mutex (since multiple goroutines might log events concurrently).
// Handle edge cases:
// Return an empty list if no events exist for a user.
// Handle file I/O errors gracefully (e.g., file not found, invalid JSON).
// Validate that action and user ID are non-empty.
// Write a main function to demonstrate logging events, retrieving user events, and saving/loading to/from a file.
// Use only Go’s standard library (e.g., encoding/json, sync, time).
// Constraints:

// User IDs and actions are non-empty strings.
// Events are stored in memory but must be persisted to a file.
// File operations should be efficient (minimize reads/writes).
// Assume the JSON file is in the same directory as the program.
// Example Behavior:

// go

// Collapse

// Wrap

// Copy
// logger := NewEventLogger()
// logger.LogEvent("login", "user123")
// logger.LogEvent("purchase", "user123")
// logger.LogEvent("logout", "user456")

// events := logger.GetUserEvents("user123") // Returns [{login, user123, timestamp1}, {purchase, user123, timestamp2}]
// err := logger.SaveToFile("events.json") // Saves events to events.json
// err = logger.LoadFromFile("events.json") // Loads events from file
