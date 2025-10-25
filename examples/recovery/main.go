package main

import (
	"fmt"
	"time"

	cheapqueue "github.com/arturwwl/cheap-queue/cheap-queue"
)

func main() {
	// This example demonstrates crash recovery
	// Run it once, kill it (Ctrl+C) before messages are consumed,
	// then run again to see messages being recovered

	var q cheapqueue.CheapQueueEngine
	q.Init("recovery-demo")

	q.Bind("persistent-queue", 20)

	// Check if there are recovered messages
	queueLen := q.QueueLen("persistent-queue")
	if queueLen > 0 {
		fmt.Printf("🔄 Recovered %d messages from previous run!\n", queueLen)
	} else {
		fmt.Println("📝 No messages to recover, publishing new ones...")
		// Publish messages
		for i := 1; i <= 5; i++ {
			msg := fmt.Sprintf("Message #%d", i)
			q.Publish("persistent-queue", []byte(msg))
			fmt.Printf("Published: %s\n", msg)
		}
	}

	fmt.Println("\n⏳ Waiting 5 seconds before consuming...")
	fmt.Println("💡 Try killing this program now (Ctrl+C) and running it again!")
	time.Sleep(5 * time.Second)

	// Consume messages
	fmt.Println("\n✅ Starting consumption...")
	q.Consume("persistent-queue", func(data []byte) {
		fmt.Printf("Consumed: %s\n", string(data))
	})

	time.Sleep(2 * time.Second)
	fmt.Println("\n✨ All messages processed!")
}
