package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

func main() {
	fmt.Println("Start")
	defer fmt.Println("Deferred 1")
	defer fmt.Println("Deferred 2")
	fmt.Println("End")
}

// the Result:
// Start
// End
// Deferred 2
// Deferred 1

// file Open and Closing

func main2() {
	file, err := os.Open("config.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close() // Always closes the file

	content, _ := ioutil.ReadAll(file)
	fmt.Println("File content:", string(content))
}

// Closing a file is related to error handling because failing to close a file can cause errors like resource leaks or file access conflicts. In your code,
// defer file.Close() ensures the file is closed even if an error occurs during reading, preventing such issues.
//  If an error happens while opening the file (err != nil), the program exits early,
//   avoiding attempts to read or close an invalid file handle, which could cause further errors.

func expensiveOperation() {
	start := time.Now()
	defer func() {
		fmt.Println("Time taken:", time.Since(start))
	}()

	// Simulate some work
	time.Sleep(2 * time.Second)
}
func main3() {
	x := 5
	defer fmt.Println("Deferred:", x) // prints 5, not 10
	x = 10
}
