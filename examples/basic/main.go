package main

import (
	"fmt"
	"time"

	cheapqueue "github.com/arturwwl/cheap-queue/cheap-queue"
)

func main() {
	// Initialize the queue engine with a unique project ID
	var q cheapqueue.CheapQueueEngine
	q.Init("basic-example")

	// Create a queue with buffer size of 10
	if err := q.Bind("tasks", 10); err != nil {
		panic(err)
	}

	// Publish some messages
	messages := []string{
		"Task 1: Process data",
		"Task 2: Send email",
		"Task 3: Update database",
	}

	for _, msg := range messages {
		if err := q.Publish("tasks", []byte(msg)); err != nil {
			fmt.Printf("Error publishing: %v\n", err)
		} else {
			fmt.Printf("Published: %s\n", msg)
		}
	}

	// Consume messages asynchronously
	q.Consume("tasks", func(data []byte) {
		fmt.Printf("Processing: %s\n", string(data))
		time.Sleep(500 * time.Millisecond) // Simulate work
	})

	// Wait for messages to be processed
	time.Sleep(3 * time.Second)

	fmt.Println("Done!")
}
