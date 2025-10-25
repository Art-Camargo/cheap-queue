package main

import "cheap-queue/engine"



func main() {
	// Application entry point
	var q engine.CheapQueueEngine
	q.Init()

	q.Bind("testQueue", 10)
	
	q.Publish("testQueue", []byte("Hello, World!"))
	q.Consume("testQueue", func(data []byte) {
		println(string(data))
	})
}
