package main

import cheapqueue "cheap-queue/cheap-queue"

func main() {

	var q cheapqueue.CheapQueueEngine
	q.Init()

	q.Bind("testQueue", 10)
	
	q.Publish("testQueue", []byte("Hello, World!"))
	q.Consume("testQueue", func(data []byte) {
		println(string(data))
	})
}
